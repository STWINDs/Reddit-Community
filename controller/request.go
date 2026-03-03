package controller

import (
	"BLUEBELL/pkg/code"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

var ErrorUserNotLogin = errors.New("用户未登录")

// getCurrentUserID 从 Gin 的上下文中获取当前用户的 ID
func GetCurrentUserID(c *gin.Context) (userID int64, err error) {
	uid, ok := c.Get(code.CtxUserIDKey)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	userID, ok = uid.(int64)
	if !ok {
		return 0, ErrorUserNotLogin
	}
	return userID, nil
}

// 将提取分页post的数据提取到公共的函数中：
func getPageInfo(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	var (
		page int64
		size int64
		err  error
	)

	page, err = strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err = strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size

}
