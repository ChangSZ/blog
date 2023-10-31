package validate

import (
	"net/http"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/gin-gonic/gin"
)

type TagStoreV struct {
}

func (tv *TagStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.TagStore
		//接收各种参数
		if err := c.ShouldBindJSON(&json); err != nil {
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &TagStore{
			Name:        json.Name,
			DisplayName: json.DisplayName,
			SeoDesc:     json.SeoDesc,
		}
		if b := appG.Validate(reqValidate); !b {
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type TagStore struct {
	Name        string `valid:"Required;MaxSize(100)"`
	DisplayName string `valid:"Required;MaxSize(100)"`
	SeoDesc     string `valid:"Required;MaxSize(250)"`
}

func (c *TagStore) Message() map[string]int {
	return map[string]int{
		"Name.Required":        403000000,
		"Name.MaxSize":         403000001,
		"DisplayName.Required": 403000002,
		"DisplayName.MaxSize":  403000003,
		"SeoDesc.Required":     403000004,
		"SeoDesc.MaxSize":      403000005,
	}
}
