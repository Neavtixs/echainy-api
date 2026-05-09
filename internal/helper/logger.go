package helper

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const RequestIDKey = "request_id"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(RequestIDKey, uuid.NewString())
		c.Next()
	}
}

func RequestLogger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		entry := NewLog(log, c).WithFields(logrus.Fields{
			"layer":      "middleware",
			"status":     c.Writer.Status(),
			"latency_ms": time.Since(start).Milliseconds(),
			"client_ip":  c.ClientIP(),
		})

		if len(c.Errors) > 0 {
			entry.WithField("errors", c.Errors.String()).Error("request completed with errors")
			return
		}

		if c.Writer.Status() >= 500 {
			entry.Error("request completed")
			return
		}

		if c.Writer.Status() >= 400 {
			entry.Warn("request completed")
			return
		}

		entry.Info("request completed")
	}
}

func NewLog(log *logrus.Logger, c *gin.Context) *logrus.Entry {
	reqID, _ := c.Get(RequestIDKey)

	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	return log.WithFields(logrus.Fields{
		"request_id": reqID,
		"path":       path,
		"method":     c.Request.Method,
	})
}
