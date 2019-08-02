package bff

import (
	"github.com/geekymedic/neon/bff/middleware"
	"github.com/geekymedic/neon/logger"

	"github.com/gin-gonic/gin"
)

var (
	_engine = gin.Default()
	_group  = _engine.Group("/api")
)

func init() {
	_group.Use(
		middleware.MetricsMiddleWare(),
		middleware.RequestTraceMiddle(logger.DefLogger(), map[string]interface{}{
			"Code":    CodeServerError,
			"Message": GetMessage(CodeServerError),
		}))
}

func Engine() *gin.RouterGroup {
	return _group
}

func MockEngine() *gin.Engine {
	return _engine
}
