package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewLog(log *logrus.Logger, c *gin.Context) *logrus.Entry {
	reqID, _ := c.Get("request_id")

	entry := log.WithFields(logrus.Fields{
		"request_id": reqID,
		"path":       c.FullPath(),
		"method":     c.Request.Method,
	})

	return entry
}
