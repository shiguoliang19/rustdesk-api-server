package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shiguoliang19/rustdesk-api-server/global"

	// "github.com/shiguoliang19/rustdesk-api-server/http/middleware"
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/shiguoliang19/rustdesk-api-server/http/router"
	"github.com/sirupsen/logrus"
)

func ApiInit() {
	gin.SetMode(global.Config.Gin.Mode)
	g := gin.New()

	//[WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
	//Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
	if global.Config.Gin.TrustProxy != "" {
		pro := strings.Split(global.Config.Gin.TrustProxy, ",")
		err := g.SetTrustedProxies(pro)
		if err != nil {
			panic(err)
		}
	}

	if global.Config.Gin.Mode == gin.ReleaseMode {
		//修改gin Recovery日志 输出为logger的输出点
		if global.Logger != nil {
			gin.DefaultErrorWriter = global.Logger.WriterLevel(logrus.ErrorLevel)
		}
	}
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 not found")
	})
	// g.Use(middleware.Logger(), gin.Recovery())

	// 自定义日志中间件
	g.Use(func(c *gin.Context) {
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
	})

	router.WebInit(g)
	router.Init(g)
	router.ApiInit(g)
	Run(g, global.Config.Gin.ApiAddr)
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
