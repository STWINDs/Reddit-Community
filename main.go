package main

import (
	"database/sql"
	"fmt"
	"os"

	// 引入你项目里的包
	"BLUEBELL/controller"
	"BLUEBELL/db" // sqlc 生成的包
	"BLUEBELL/logger"
	"BLUEBELL/logic"
	"BLUEBELL/pkg/snowflake"
	"BLUEBELL/redis"
	"BLUEBELL/router"
	"BLUEBELL/setting"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// initDB 初始化并返回数据库连接对象
// 注意：不要在这里 defer db.Close()，否则连接会立即关闭
func initDB() (*sql.DB, error) {
	// 如果你用了 godotenv，可以在这里加载，也可以在 setting 里统一加载
	godotenv.Load(".env")

	// 这里建议用 setting.Conf 里的配置，如果没有就用环境变量
	// dbstr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local", ...)
	// 暂时沿用你原本的写法：
	dbstr := os.Getenv("DATABASE_URL")

	conn, err := sql.Open("mysql", dbstr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// @title Bluebell Forum API
// @version 1.0
// @description Bluebell 论坛后端接口文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

func main() {
	// 1. 加载配置
	if err := setting.Init("config.yaml"); err != nil {
		fmt.Printf("load config failed, err:%v\n", err)
		return
	}

	// 2. 初始化日志
	if err := logger.Init(setting.Conf.LogConfig, setting.Conf.Mode); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	defer zap.L().Sync() // 确保日志刷入磁盘

	// 3. 初始化翻译器 (Validator)
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("init trans failed, err:%v\n", err)
		return
	}

	// 3.5. 初始化 Snowflake ID 生成器
	if err := snowflake.InitSnowflake(setting.Conf.StartTime, setting.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err:%v\n", err)
		return
	}

	// 3.6. 允许环境变量覆盖 Redis 配置 (适配 Docker)
	if host := os.Getenv("REDIS_HOST"); host != "" {
		setting.Conf.RedisConfig.Host = host
	}
	if port := os.Getenv("REDIS_PORT"); port != "" {
		fmt.Sscanf(port, "%d", &setting.Conf.RedisConfig.Port)
	}

	// 4. 【核心】初始化数据库连接&redis
	conn, err := initDB()
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	defer conn.Close() // 程序退出时关闭数据库连接

	if err := redis.Init(setting.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err:%v\n", err)
		return
	}
	defer redis.Close()

	// ================= 依赖注入开始 =================

	// 5. 初始化 Data Access Layer (sqlc)
	// 这里的 queries 就是 "电钻"
	queries := db.New(conn)

	// 6. 初始化 Logic Layer
	// 把 "电钻" 交给 "包工头" (UserLogic)
	userLogic := logic.NewUserLogic(queries)
	communityLogic := logic.NewCommunityLogic(queries)
	postLogic := logic.NewPostLogic(queries)

	// 7. 初始化 Controller Layer
	// 把 "包工头" 介绍给 "接待员" (UserHandler)
	userHandler := controller.NewUserHandler(userLogic)
	communityHandler := controller.NewCommunityHandler(communityLogic)
	posthandler := controller.NewPostHandler(postLogic)

	// 8. 初始化 Router
	// 把 "接待员" 带到 "大厅" (Router)
	r := router.SetupRouter(setting.Conf.Mode, userHandler, communityHandler, posthandler)

	// ================= 依赖注入结束 =================

	// 9. 启动服务
	// 注意：r.Run() 是阻塞的，后面的代码不会执行，所以不要在它后面写 r.POST
	err = r.Run(fmt.Sprintf(":%d", setting.Conf.Port))
	if err != nil {
		fmt.Printf("run server failed, err:%v\n", err)
		return
	}
}
