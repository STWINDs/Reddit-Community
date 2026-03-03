package controller

import (
	"BLUEBELL/logic"
	"BLUEBELL/models"
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PostHandler struct {
	logic *logic.PostLogic // 持有对应的社区逻辑层
}

func NewPostHandler(l *logic.PostLogic) *PostHandler {
	return &PostHandler{logic: l}
}

func (h *PostHandler) CreatePostHandler(c *gin.Context) {
	// 1. 参数校验

	// c.ShouldBindJSON() 绑定 JSON 参数到结构体，并进行校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Debug("c.ShouldBindJSON failed", zap.Any("err", err))
		zap.L().Error("CreatePostHandler with invalid param", zap.Error(err))
		response.RespondWithError(c, code.CodeInvalidParam) // 这里我们也返回 200，但 message 里有错误信息，保持和登录接口一致的风格
		return
	}

	//从上下文中获取用户ID
	userID, err := GetCurrentUserID(c)
	if err != nil {
		zap.L().Error("userID not found in context")
		response.RespondWithError(c, code.CodeServerBusy)
		return
	}
	p.AuthorID = userID // 将用户ID设置为帖子作者ID
	// 2. 业务处理 (调用 Logic 层)

	if err := h.logic.CreatePost(c.Request.Context(), p); err != nil {
		zap.L().Error("logic.CreatePost failed", zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy)

	}
	// 3. 返回成功响应
	response.RespondWithSuccess(c, nil)
}

// 获取帖子详情的handler
func (l *PostHandler) GetPostDetailHandler(c *gin.Context) {
	// 1. 从 URL 中获取 post_id 参数
	// Param() 方法获取 URL 中的参数，参数名是路由中定义的占位符名称，这里我们假设路由定义为 "/post/:id"，所以占位符名称是 "id"
	// 如果取不到或则为空
	postIDStr := c.Param("id")
	if postIDStr == "" {
		zap.L().Error("post_id is required")
		response.RespondWithError(c, code.CodeInvalidParam)
		return
	}

	// 2. 根据id查询数据库，获取帖子详情
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		zap.L().Error("invalid post_id", zap.String("post_id", postIDStr), zap.Error(err))
		response.RespondWithError(c, code.CodeInvalidParam)
		return
	}

	data, err := l.logic.GetPostDetail(c.Request.Context(), postID)
	if err != nil {
		zap.L().Error("logic.GetPostDetail() failed", zap.Int64("post_id", postID), zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy)
		return
	}
	// 3. 返回响应
	response.RespondWithSuccess(c, data)

}

//获取post列表

func (l *PostHandler) GetPostListHandler(c *gin.Context) {
	// 获取数据
	// 获取分页参数
	page, size := getPageInfo(c)
	data, err := l.logic.GetPostList(c, page, size)
	if err != nil {
		zap.L().Error("l.logic.GetPostList() failed", zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy)
		return
	}

	response.RespondWithSuccess(c, data)

}

// 根据时间或分数获取帖子列表接口
// 根据前端传来的参数动态获取帖子列表（按create_time排序或按照score排序）
// 1.获取参数
// 2.去redis查询id值
// 3.根据id去数据库查询帖子详情信息

func (l *PostHandler) GetPostListHandler2(c *gin.Context) {
	// 获取数据
	// 获取分页参数
	p := &models.PostListParam{}
	p.Page = 1
	p.Size = 10
	p.Order = models.OrderTime

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetPostListHandler2 with Invalid Params", zap.Error(err))
		response.RespondWithError(c, code.CodeInvalidParam)
		return
	}

	//GET（query string参数）: /api/v1/posts2?page=1&size=10&order=time
	data, err := l.logic.GetPostList2(c, p)
	if err != nil {
		zap.L().Error("l.logic.GetPostList2() failed", zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy)
		return
	}

	response.RespondWithSuccess(c, data)

}

// 根据community 来查询帖子列表
func (l *PostHandler) GetCommunityPostListHandler(c *gin.Context) {

	p := &models.CommunityPostListParam{}
	p.Page = 1
	p.Size = 10
	p.Order = models.OrderTime

	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("GetCommunityPostListHandler with Invalid Params", zap.Error(err))
		response.RespondWithError(c, code.CodeInvalidParam)
		return
	}

	data, err := l.logic.GetCommunityPostList(c, p)
	if err != nil {
		zap.L().Error("l.logic.GetCommunityPostList() failed", zap.Error(err))
		response.RespondWithError(c, code.CodeServerBusy)
		return
	}

	response.RespondWithSuccess(c, data)

}
