package extend

import (
	"github.com/geekymedic/neon/logger"
	"testing"
)

func TestNewSimpleLog(t *testing.T) {
	sLog := NewSimpleLog()
	sLog.Warn("Hello")
	sLog.With("method", "get", "path", "/api/admin", "arg", "").Warn("trace log")
	logger.SetLogger(sLog)

	logger.Info("ddd")
}
