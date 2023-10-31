package middleware

import (
	"net/http"
	"strings"

	"github.com/ChangSZ/blog/infra/log"
	"github.com/gin-gonic/gin"
)

func CheckExist() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		s := strings.Contains(path, "/backend/")
		c.Next()
		status := c.Writer.Status()
		if status == 404 {
			log.WithTrace(c).Errorf("routerNoFound CheckExist, Path: %v", path)
			if s {
				c.Redirect(http.StatusMovedPermanently, "/backend/")
			} else {
				c.Redirect(http.StatusMovedPermanently, "/404")
			}
		}
	}
}
