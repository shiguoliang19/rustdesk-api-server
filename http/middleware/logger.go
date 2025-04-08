package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shiguoliang19/rustdesk-api-server/global"
	"github.com/sirupsen/logrus"

	"bytes"
	"fmt"
	"io"
	"log"
	"time"
)

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		global.Logger.WithFields(
			logrus.Fields{
				"uri":    c.Request.URL,
				"ip":     c.ClientIP(),
				"method": c.Request.Method,
			}).Debug("Request")
		c.Next()
	}
}

func Logger2() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 你的自定义格式
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Logger3() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ip := c.ClientIP()
		headers := c.Request.Header

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
		}

		// 创建响应写入器
		writer := &responseWriter{ResponseWriter: c.Writer, body: bytes.NewBuffer(nil)}
		c.Writer = writer

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		responseBody := writer.body.String()

		log.Printf("Request: %s %s | IP: %s | Status: %d | Duration: %v | Headers: %v | RequestBody: %s | ResponseBody: %s",
			c.Request.Method,
			c.Request.URL.Path,
			ip,
			status,
			duration,
			headers,
			string(requestBody),
			responseBody,
		)
	}
}
