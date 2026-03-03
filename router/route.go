package router

import (
	"BLUEBELL/controller"
	"BLUEBELL/logger"
	"BLUEBELL/middlewares"
	"BLUEBELL/pkg/code"
	"BLUEBELL/pkg/response"
	"net/http"
	_ "net/http/pprof"
	"time"

	_ "BLUEBELL/docs"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 接收 userHandler 实例并返回配置好的 gin 引擎
func SetupRouter(mode string, userHandler *controller.UserHandler, communityHandler *controller.CommunityHandler, postHandler *controller.PostHandler) *gin.Engine {
	// 设置 Gin 模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// 使用自定义限流中间件，每 500 毫秒放一个令牌，桶容量为 100
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middlewares.RateLimitMiddleware(500*time.Millisecond, 1000))
	// r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	// 注册业务路由
	v1.POST("/signup", userHandler.SignUpHandler)

	// 注意：这里注册登录路由，调用 userHandler 的 LoginHandler 方法
	// 这样就完成了从路由到控制器的连接，用户访问 /login 路径时，就会调用 userHandler.LoginHandler 方法来处理登录请求。
	v1.POST("/login", userHandler.LoginHandler)

	// 调试路由：在非 release 模式下允许不带 JWT 请求帖子详情，方便本地调试
	if mode != gin.ReleaseMode {
		v1.GET("/debug/post/:id", postHandler.GetPostDetailHandler)
	}

	v1.Use(middlewares.JWTAuthMiddleware())
	{
		v1.GET("/community", communityHandler.GetCommunityHandler)
		v1.GET("/community/:id", communityHandler.GetCommunityDetailHandler)

		v1.POST("/post", postHandler.CreatePostHandler)
		v1.GET("/post/:id", postHandler.GetPostDetailHandler)
		v1.GET("/posts", postHandler.GetPostListHandler)
		//根据时间/分数获得帖子列表
		v1.GET("/posts2", postHandler.GetPostListHandler2)
		v1.POST("/vote", postHandler.PostVoteController)
	}
	// 需要认证的路由组，后续的路由都需要 JWT 认证

	// 系统检测路由
	v1.GET("/ping", middlewares.JWTAuthMiddleware(), func(c *gin.Context) {
		// 这里我们可以通过之前设置在 JWTAuthMiddleware 中的用户 ID 来验证是否正确获取了用户信息
		// 1. 先调用函数，接收两个返回值
		userID, err := controller.GetCurrentUserID(c)

		// 2. 检查错误
		if err != nil {
			// 如果获取不到 userID，说明没登录，直接返回错误响应
			response.RespondWithError(c, code.CodeNeedLogin)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
			"user_id": userID, // 这里我们可以通过之前设置在 JWTAuthMiddleware 中的用户 ID 来验证是否正确获取了用户信息
		})
	})

	pprof.Register(r)
	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "404 not found"})
	})

	return r
}
