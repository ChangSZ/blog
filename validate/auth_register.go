package validate

import (
	"net/http"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/gin-gonic/gin"
)

type AuthRegisterV struct {
}

func (av *AuthRegisterV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.AuthRegister
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &AuthRegister{
			Email:    json.Email,
			Password: json.Password,
			UserName: json.UserName,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type AuthRegister struct {
	UserName string `valid:"Required;MaxSize(30)"`
	Email    string `valid:"Required;Email"`
	Password string `valid:"Required;MaxSize(30)"`
}

func (av *AuthRegister) Message() map[string]int {
	return map[string]int{
		"Email.Required":    407000000,
		"Email.Email":       407000001,
		"Password.Required": 407000002,
		"Password.MaxSize":  407000003,
		"UserName.Required": 407000012,
		"UserName.MaxSize":  407000013,
	}
}
