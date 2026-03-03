package logic

import (
	"BLUEBELL/models"
	"BLUEBELL/redis"
	"strconv"

	"go.uber.org/zap"
)

//每个帖子只允许发布后一周内投票
//1.到期后将redis保存的zset 赞成票数和反对票数保存到mysql表中
//2.到期后删除KeyPostVotedZsetPF

func (h *PostLogic) VoteForPost(userID int64, p *models.VoteDataParam) error {
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.Int64("postID", p.PostID),
		zap.Int8("direction", p.Direction))

	// 假设你的 redis.VoteForPost 接收 string 类型的 ID 进行拼接 Key
	// 我们手动转换一下
	strUserID := strconv.FormatInt(userID, 10)
	strPostID := strconv.FormatInt(p.PostID, 10)
	return redis.VoteForPost(strUserID, strPostID, float64(p.Direction))
}
