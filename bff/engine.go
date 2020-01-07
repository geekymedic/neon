package bff

import (
	"github.com/gin-gonic/gin"
)

var (
	_engine = gin.Default()
	_group  = _engine.Group("/api")
)

func init() {
	_group.Use(
		MetricsMiddleWare(),
		RequestTraceMiddle(map[string]interface{}{
			"Code":    CodeServerError,
			"Message": GetMessage(CodeServerError),
		}))
}

func Engine() *gin.RouterGroup {
	return _group
}

func RootEngine() *gin.Engine {
	return _engine
}

func MockEngine() *gin.Engine {
	return _engine
}
