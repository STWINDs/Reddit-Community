package logic

import (
	"BLUEBELL/db"
	"BLUEBELL/models"
	"BLUEBELL/pkg/snowflake"
	"BLUEBELL/redis"
	"context"
	"database/sql"
	"strings"

	"go.uber.org/zap"
)

type PostLogic struct {
	store *db.Queries // 依赖注入
}

func NewPostLogic(store *db.Queries) *PostLogic {
	return &PostLogic{store: store}
}

func (h *PostLogic) CreatePost(ctx context.Context, p *models.Post) error {
	// 生成post_id

	p.ID = int64(snowflake.GenID())

	// 保存到数据库
	err := h.store.CreatePost(ctx, db.CreatePostParams{
		PostID:      uint64(p.ID),
		AuthorID:    uint64(p.AuthorID),
		CommunityID: uint64(p.CommunityID),
		Title:       p.Title,
		Content:     sql.NullString{String: p.Content, Valid: true},
	})
	if err != nil {
		zap.L().Error("db.CreatePost() failed", zap.Error(err))
		return err
	}

	// 将贴子ID添加到Redis的Sorted Set中
	if err := redis.AddPostIDToRedis(p.ID, p.CommunityID); err != nil {
		zap.L().Error("redis.AddPostIDToRedis() failed", zap.Error(err))
		// 注意：这里不返回错误，因为即使Redis写入失败，贴子已经保存到数据库了
		// 可以选择返回错误，取决于业务需求
	}

	// 返回结果
	return nil
}

// 获取帖子详情的逻辑函数
func (h *PostLogic) GetPostDetail(ctx context.Context, postID int64) (data *models.ApiPostDetail, err error) {
	// 从数据库查询帖子详情

	post, err := h.store.GetPostDetailByID(ctx, uint64(postID))
	if err != nil {
		zap.L().Error("db.GetPostDetailByID() failed", zap.Int64("postID", postID), zap.Error(err))
		return
	}

	// 查作者信息
	// Use new GetUserByUserID to lookup by snowflake user_id instead of internal id
	user, err := h.store.GetUserByUserID(ctx, int64(post.AuthorID))
	if err != nil {
		zap.L().Error("h.store.GetUserByUserID() failed", zap.Int64("author_id", int64(post.AuthorID)), zap.Error(err))
		return
	}

	// 查社区信息
	community, err := h.store.GetCommunityDetailByID(ctx, post.CommunityID)
	if err != nil {
		zap.L().Error("h.store.GetCommunityDetailByID() failed", zap.Int64("community_id", int64(post.CommunityID)), zap.Error(err))
		return
	}

	data = &models.ApiPostDetail{
		AuthorName:                user.Username,
		GetPostDetailByIDRow:      &post,
		GetCommunityDetailByIDRow: &community,
	}

	// 返回帖子详情
	return

}

// 获取帖子列表 (使用 SQL JOIN 查询，高效无 N+1)
// 分页获取Post
func (h *PostLogic) GetPostList(ctx context.Context, page int64, size int64) (data []*models.ApiPostListItem, err error) {

	offset := (page - 1) * size

	// 2. 构造 sqlc 需要的 Params 结构体
	params := db.GetPostListWithDetailsParams{
		Offset: int32(offset),
		Limit:  int32(size),
	}
	posts, err := h.store.GetPostListWithDetails(ctx, params)
	if err != nil {
		zap.L().Error("db.GetPostListWithDetails() failed", zap.Error(err))
		return nil, err
	}

	data = make([]*models.ApiPostListItem, 0, len(posts))

	for _, post := range posts {
		item := &models.ApiPostListItem{
			ID:            post.ID,
			AuthorID:      post.AuthorID,
			AuthorName:    post.AuthorName.String,
			CommunityID:   post.CommunityID,
			CommunityName: post.CommunityName.String,
			Title:         post.Title,
			Content:       post.Content.String,
			Status:        post.Status,
			CreateTime:    post.CreateTime,
		}
		data = append(data, item)
	}

	return data, nil
}

// 2.去redis查询id值

func (h *PostLogic) GetPostList2(ctx context.Context, p *models.PostListParam) (data []*models.ApiPostListItem, err error) {
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 1. 将 []string {"10", "1", "5"} 转换为 "10,1,5"
	idStr := strings.Join(ids, ",")

	// 2. 构造参数

	// 3. 调用函数
	// 根据id去MySQL数据库查询帖子详情信息
	posts, err := h.store.GetPostListByIDs(ctx, db.GetPostListByIDsParams{
		FINDINSET:   idStr,
		FINDINSET_2: idStr,
	})
	if err != nil {
		zap.L().Error("h.store.GetPostListByIDs failed", zap.Error(err))
		return nil, err
	}
	//提前查询号每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 数据聚合：补充作者和社区信息
	data = make([]*models.ApiPostListItem, 0, len(posts))

	for idx, post := range posts {
		// A. 查询作者信息 (实际项目中建议这里也做批量查询或加缓存)
		user, err := h.store.GetUserByUserID(ctx, int64(post.AuthorID))
		authorName := "已注销用户"
		if err == nil {
			authorName = user.Username
		}

		// B. 查询社区信息
		community, err := h.store.GetCommunityDetailByID(ctx, post.CommunityID)
		communityName := "未知社区"
		if err == nil {
			communityName = community.Name
		}

		// C. 拼装 ApiPostListItem
		item := &models.ApiPostListItem{
			ID:            post.PostID,
			AuthorID:      post.AuthorID,
			AuthorName:    authorName,
			CommunityID:   post.CommunityID,
			CommunityName: communityName,
			Title:         post.Title,
			// 处理 sql.NullString: 如果 Valid 为 true 则取值，否则为空字符串
			Content:    post.Content.String,
			Status:     1, // 默认正常
			CreateTime: post.CreateTime,
			VoteNum:    voteData[idx],
		}
		data = append(data, item)
	}

	return data, nil

}

func (h *PostLogic) GetCommunityPostList(ctx context.Context, p *models.CommunityPostListParam) (data []*models.ApiPostListItem, err error) {

	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}

	if len(ids) == 0 {
		zap.L().Warn("redis.GetPostIDsInOrder(p) return 0 data")
		return
	}
	// 1. 将 []string {"10", "1", "5"} 转换为 "10,1,5"
	idStr := strings.Join(ids, ",")

	// 2. 构造参数

	// 3. 调用函数
	// 根据id去MySQL数据库查询帖子详情信息
	posts, err := h.store.GetPostListByIDs(ctx, db.GetPostListByIDsParams{
		FINDINSET:   idStr,
		FINDINSET_2: idStr,
	})
	if err != nil {
		zap.L().Error("h.store.GetPostListByIDs failed", zap.Error(err))
		return nil, err
	}
	//提前查询号每篇帖子的投票数
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}

	// 3. 数据聚合：补充作者和社区信息
	data = make([]*models.ApiPostListItem, 0, len(posts))

	for idx, post := range posts {
		// A. 查询作者信息 (实际项目中建议这里也做批量查询或加缓存)
		user, err := h.store.GetUserByUserID(ctx, int64(post.AuthorID))
		authorName := "已注销用户"
		if err == nil {
			authorName = user.Username
		}

		// B. 查询社区信息
		community, err := h.store.GetCommunityDetailByID(ctx, post.CommunityID)
		communityName := "未知社区"
		if err == nil {
			communityName = community.Name
		}

		// C. 拼装 ApiPostListItem
		item := &models.ApiPostListItem{
			ID:            post.PostID,
			AuthorID:      post.AuthorID,
			AuthorName:    authorName,
			CommunityID:   post.CommunityID,
			CommunityName: communityName,
			Title:         post.Title,
			// 处理 sql.NullString: 如果 Valid 为 true 则取值，否则为空字符串
			Content:    post.Content.String,
			Status:     1, // 默认正常
			CreateTime: post.CreateTime,
			VoteNum:    voteData[idx],
		}
		data = append(data, item)
	}

	return data, nil

}
