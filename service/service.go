package service

import (
	"errors"
	"net/http"
	"time"

	"github.com/NYTimes/gziphandler"
	"github.com/gin-gonic/gin"
	"github.com/guomio/go-template/logger"
	"github.com/guomio/go-template/tools"
	"github.com/rs/cors"
)

var (
	// g gin 实例
	g *gin.Engine

	// StatusBadCode 错误代码
	StatusBadCode = "400"
	// StatusOkCode 正常代码
	StatusOkCode = "00"

	// ErrForbidden forbidden error
	ErrForbidden = errors.New("forbidden")
	// ErrBindings invalid param bindings
	ErrBindings = errors.New("params are required")
)

// Init 初始化Gin实例
func Init() {
	gin.SetMode(gin.ReleaseMode)
	g = gin.New()
	g.Use(DefaultLogger(), gin.Recovery())

	// 此处注册路由
}

// Group gin Group
func Group(relativePath ...string) *gin.RouterGroup {
	return g.Group(tools.CombineURLs(relativePath...))
}

// Handler http handler
func Handler(opt cors.Options) http.Handler {
	c := cors.New(opt)
	return gziphandler.GzipHandler(c.Handler(g))
}

// DefaultCorsOption 默认跨域配置
func DefaultCorsOption() cors.Options {
	return cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}
}

// Response 服务器返回信息
type Response struct {
	Data   interface{} `json:"data"`
	Status string      `json:"status"`
	Msg    string      `json:"msg"`
	Error  string      `json:"error"`
}

// ErrorResponse 错误返回
func ErrorResponse(c *gin.Context, err, msg error) {
	c.JSON(http.StatusOK, Response{Data: nil, Status: StatusBadCode, Error: err.Error(), Msg: msg.Error()})
}

// SuccessResponse 成功返回
func SuccessResponse(c *gin.Context, d interface{}) {
	c.JSON(http.StatusOK, Response{Data: d, Status: StatusOkCode})
}

//HandleResponse response handler
func HandleResponse(c *gin.Context, d interface{}, err, msg error) {
	if err != nil {
		ErrorResponse(c, err, msg)
	} else {
		SuccessResponse(c, d)
	}
}

// DefaultLogger 默认日志
func DefaultLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if param.Latency > time.Minute {
			param.Latency = param.Latency - param.Latency%time.Second
		}

		go func() {
			logger.L.Info("[MIO]",
				logger.Int("status", param.StatusCode),
				logger.String("method", param.Method),
				logger.String("duration", param.Latency.String()),
				logger.String("ip", param.ClientIP),
				logger.String("path", param.Path),
			)
		}()
		return ""
	})
}
