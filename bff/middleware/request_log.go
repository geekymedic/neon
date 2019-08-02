package middleware

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/geekymedic/neon"
	"github.com/geekymedic/neon/bff/types"
	"github.com/geekymedic/neon/logger"
	"github.com/geekymedic/neon/version"

	"github.com/gin-gonic/gin"
)

func RequestTraceMiddle(log logger.Logger, failOut map[string]interface{}, ignore ...string) gin.HandlerFunc {
	var skips = map[string]struct{}{}
	for _, _ignore := range ignore {
		skips[_ignore] = struct{}{}
	}
	return func(c *gin.Context) {

		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		session := neon.NewSessionFromGinCtx(c)
		defer func() {
			err := recover()

			// Log only when path is not being skipped
			if _, ok := skips[path]; !ok {
				param := gin.LogFormatterParams{
					Request: c.Request,
					Keys:    c.Keys,
				}

				// Stop timer
				param.TimeStamp = time.Now()
				param.Latency = param.TimeStamp.Sub(start)

				param.ClientIP = c.ClientIP()
				param.Method = c.Request.Method
				param.StatusCode = c.Writer.Status()
				param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

				param.BodySize = c.Writer.Size()

				if raw != "" {
					path = path + "?" + raw
				}
				param.Path = path

				var body []byte
				var err error
				if cb, ok := c.Get(gin.BodyBytesKey); ok {
					if cbb, ok := cb.([]byte); ok {
						body = cbb
					}
				}
				if body == nil {
					body, err = ioutil.ReadAll(c.Request.Body)
					if err != nil {
						// TODO return it
					} else {
						c.Set(gin.BodyBytesKey, body)
					}
				}
				bodySize := len(body)
				if bodySize > 1<<10 {
					bodySize = 1 << 10
				}

				contentSize := c.Request.Header.Get("Content-Length")
				sessionLog := session.ShortLog()
				log := log.With("pro_name", version.PRONAME, "gitcommit", version.GITCOMMIT, "path", param.Path,
					"method", param.Method,
					"status", param.StatusCode,
					"req_size", contentSize,
					"resp_size", param.BodySize,
					"latency", fmt.Sprintf("%v", param.Latency),
					"client_ip", param.ClientIP,
				).With(sessionLog...)
				if param.ErrorMessage != "" {
					log.With("err", param.ErrorMessage).Error("http request trace")
					log.With("trace", session.Trace,
						"body", string(body)[:bodySize]).Error("request body")
				} else {
					log.Info("http request trace")
					log.With("trace", session.Trace,
						"body", string(body)[:bodySize]).Info("request body")
				}
			}
			if err != nil {
				if !c.Writer.Written() {
					c.JSON(http.StatusOK, failOut)
					c.Set(types.ResponseStatusCode, failOut["Code"])
				}

				log.With("panic", err).Error("process panic")
				var buf [1024]byte
				runtime.Stack(buf[:], true)
				log.With("stack", fmt.Sprintf("%s", buf)).Error("panic stack")
			}
		}()

		// Process request
		c.Next()
	}
}
