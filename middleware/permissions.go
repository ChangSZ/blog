package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/infra/jwt"
)

func Permission(routerAsName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiG := api.Gin{C: c}
		fmt.Println(routerAsName, c.Request.Method)
		res := common.CheckPermissions(routerAsName, c.Request.Method)
		if !res {
			log.WithTrace(c).Error("router permission")
			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}

		token := c.GetHeader("x-auth-token")
		if routerAsName == "console.post.imgUpload" {
			token = c.GetHeader("X-Upload-Token")
		}

		if token == "" {
			log.WithTrace(c).Error("token null")
			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}

		userId, err := jwt.ParseToken(c, token)
		if err != nil {
			log.WithTrace(c).Errorf("parse token error, err: %v", err)
			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			log.WithTrace(c).Errorf("strconv token error, err: %v", err)
			apiG.Response(http.StatusOK, 400001005, nil)
			return
		}
		c.Set("userId", userIdInt)
		c.Set("token", token)
		c.Next()
	}
}
