package models

import (
	"BLUEBELL/db"
)

const (
	OrderTime  = `Time`
	Orderscore = `score`
)

// SignupParams 注册参数
type SignupParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// LoginParams 登录参数
type LoginParams struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Token    string
}

type Post struct {
	ID          int64  `json:"id,string"`
	AuthorID    int64  `json:"author_id"`
	CommunityID int64  `json:"community_id" binding:"required"`
	Status      int64  `json:"status"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content" binding:"required"`
	CreateTime  string `json:"create_time"`
}

// 帖子详情接口的结构体
type ApiPostDetail struct {
	AuthorName string `json:"author_name"`
	*db.GetPostDetailByIDRow
	*db.GetCommunityDetailByIDRow `json:"community"`
}

// 帖子列表项目结构体 (对应 db.GetPostListWithDetailsRow)
type ApiPostListItem struct {
	ID            uint64      `json:"id"`
	AuthorID      uint64      `json:"author_id"`
	AuthorName    string      `json:"author_name"`
	CommunityID   uint64      `json:"community_id"`
	VoteNum       int64       `json:"vote_num"`
	CommunityName string      `json:"community_name"`
	Title         string      `json:"title"`
	Content       string      `json:"content"`
	Status        int8        `json:"status"`
	CreateTime    interface{} `json:"create_time"`
}

type VoteDataParam struct {
	// UserID从request中获取
	PostID    int64 `json:"post_id,string" binding:"required"`
	Direction int8  `json:"direction,string" binding:"required,oneof= 1 0 -1"` //赞同1，反对-1，取消0
}

// 获取帖子列表的Query String参数,这个是url形式的数据
type PostListParam struct {
	Page  int64  `json:"page" form:"page"`
	Size  int64  `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}

type CommunityPostListParam struct {
	PostListParam
	CommunityID int64 `json:"community_id" form:"community_id"`
}
