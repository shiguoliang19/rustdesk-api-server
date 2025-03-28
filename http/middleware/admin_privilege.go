package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/shiguoliang19/rustdesk-api-server/http/response"
	"github.com/shiguoliang19/rustdesk-api-server/service"
)

// AdminPrivilege ...
func AdminPrivilege() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := service.AllService.UserService.CurUser(c)

		if !service.AllService.UserService.IsAdmin(u) {
			response.Fail(c, 403, "无权限")
			c.Abort()
			return
		}

		c.Next()
	}
}
