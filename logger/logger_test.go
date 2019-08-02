package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	Info("Hello Word")
	Warn("Hello Word")
	log := With("service", "user_login")
	log.Info("hello word")
}

func TestWarn(t *testing.T) {
	log := logrus.New()
	var entry = log.WithFields(map[string]interface{}{"service": "login"})
	entry.Logger.Info("login")
}

func TestSetLevel(t *testing.T) {
	os.Setenv("neon_LOG_OUT", "json")
	setJson()
	t.Run("logrus", func(t *testing.T) {
		entry := logrus.New()
		entry.SetLevel(DebugLevel)
		entry.Debug("dddddd")
	})
	t.Run("inner", func(t *testing.T) {
		log.SetLevel(DebugLevel)
		//SetLevel(DebugLevel)
		log.Debug("Debug ...")
		log.SetLevel(InfoLevel)
		log.Debug("Info ...")
		Info("ddddd")
	})
}