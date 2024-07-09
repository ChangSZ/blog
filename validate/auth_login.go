package validate

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/infra/gin/api"
)

type AuthLoginV struct {
}

func (av *AuthLoginV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.AuthLogin
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &AuthLogin{
			Email:      json.Email,
			Password:   json.Password,
			Captcha:    json.Captcha,
			CaptchaKey: json.CaptchaKey,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type AuthLogin struct {
	Email      string `valid:"Required;Email"`
	Password   string `valid:"Required;MaxSize(30)"`
	Captcha    string `valid:"Required;MaxSize(5)"`
	CaptchaKey string `valid:"Required;MaxSize(30)"`
}

func (av *AuthLogin) Message() map[string]int {
	return map[string]int{
		"Email.Required":      407000000,
		"Email.Email":         407000001,
		"Password.Required":   407000002,
		"Password.MaxSize":    407000003,
		"Captcha.Required":    407000004,
		"Captcha.MaxSize":     407000005,
		"CaptchaKey.Required": 407000006,
		"CaptchaKey.MaxSize":  407000007,
	}
}
