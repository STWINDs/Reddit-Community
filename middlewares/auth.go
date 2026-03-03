package middlewares

import (
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/jwt"
	"BLUEBELL/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			// 没有携带 Authorization，返回需要登录的错误信息
			response.RespondWithError(c, code.CodeNeedLogin)
			c.Abort()
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.RespondWithError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.RespondWithError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}

		if mc.UserID == 0 {
			// 提示 Token 非法（不要告诉黑客具体原因）
			response.RespondWithError(c, code.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		//Set方法可以在上下文中存储一个键值对，这样在后续的处理函数中就可以通过Get方法来获取这个值了。
		c.Set(code.CtxUserIDKey, mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get(CtxUserIDKey)来获取当前请求的用户信息
	}
}
