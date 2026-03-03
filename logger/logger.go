package logger

import (
	"BLUEBELL/setting"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lg *zap.Logger

func Init(cfg *setting.LogConfig, mode string) (err error) {
	//getLogWriter 函数负责创建一个 zapcore.WriteSyncer，用于将日志写入到指定的文件中。它使用了 lumberjack 包来实现日志文件的滚动功能，确保日志文件不会过大。lumberjack.Logger 结构体包含了日志文件的名称、最大大小、最大备份数量和最大年龄等配置项。
	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	// 获取日志编码器，根据需要可以选择 JSONEncoder 或 ConsoleEncoder，这里我们使用 JSONEncoder 以便于日志的结构化输出，方便后续的日志分析和处理。
	encoder := getEncoder()
	// 解析日志级别，zapcore.Level 是 zap 包中定义的日志级别类型，UnmarshalText 方法可以将字符串形式的日志级别转换为 zapcore.Level 类型的值。这样我们就可以在配置文件中使用字符串来指定日志级别，例如 "debug"、"info"、"warn"、"error" 等。
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return
	}
	var core zapcore.Core
	//mode 是一个字符串变量，表示当前的运行模式，可以是 "dev"（开发模式）或 "prod"（生产模式）。根据不同的运行模式，我们可以配置不同的日志输出方式。在开发模式下，我们希望日志能够同时输出到文件和终端，以便于开发者在调试过程中能够实时看到日志信息；而在生产模式下，我们通常只需要将日志输出到文件中，以便于后续的分析和处理。
	if mode == "dev" {
		// 进入开发模式，日志输出到终端
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee(
			zapcore.NewCore(encoder, writeSyncer, l),
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
		//NewTee 可以同时输出到多个目的地，这里我们既输出到文件，也输出到终端，终端的日志级别设置为 DebugLevel，确保开发过程中能看到所有日志信息。
		// 这种方式在开发环境中非常有用，可以让开发者在终端实时看到日志输出，同时又能将日志持久化到文件中，方便后续分析和调试。
		// NewCore 是 zap 的核心组件，用于创建一个日志核心，参数分别是编码器、日志输出目的地和日志级别。在这里，我们为文件输出创建了一个核心，并为终端输出创建了另一个核心。
		// Lock 是 zapcore 包中的一个函数，用于确保在多线程环境下对 os.Stdout 的写入是安全的。它会返回一个 WriteSyncer，确保多个 goroutine 同时写入日志时不会发生竞态条件。
		// writeSyncer 是一个 zapcore.WriteSyncer 接口的实现，负责将日志写入到指定的目的地。在这里，我们使用了 lumberjack 包来实现日志文件的滚动功能，确保日志文件不会过大。
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}

	lg = zap.New(core, zap.AddCaller())

	zap.ReplaceGlobals(lg)
	zap.L().Info("init logger success")
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		lg.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					lg.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					lg.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
