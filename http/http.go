package http

import (
	"github.com/gin-gonic/gin"
	"github.com/shiguoliang19/rustdesk-api-server/global"

	"github.com/shiguoliang19/rustdesk-api-server/http/middleware"
	"net/http"
	"strings"

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
	g.Use(middleware.Logger3(), gin.Recovery())

	router.WebInit(g)
	router.Init(g)
	router.ApiInit(g)
	Run(g, global.Config.Gin.ApiAddr)
}
