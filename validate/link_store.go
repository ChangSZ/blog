package validate

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/infra/gin/api"
)

type LinkStoreV struct {
}

func (lv *LinkStoreV) MyValidate() gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := api.Gin{C: c}
		var json common.LinkStore
		if err := c.ShouldBindJSON(&json); err != nil {
			fmt.Println(json)
			appG.Response(http.StatusOK, 400001000, nil)
			return
		}

		reqValidate := &LinkStore{
			Name:  json.Name,
			Link:  json.Link,
			Order: json.Order,
		}
		if b := appG.Validate(reqValidate); !b {
			fmt.Println(reqValidate, json)
			return
		}
		c.Set("json", json)
		c.Next()
	}
}

type LinkStore struct {
	Name  string `valid:"Required;MaxSize(100)"`
	Link  string `valid:"Required;MaxSize(100)"`
	Order int    `valid:"Min(0)"`
}

func (c *LinkStore) Message() map[string]int {
	return map[string]int{
		"Name.Required": 406000000,
		"Name.MaxSize":  406000001,
		"Link.Required": 406000002,
		"Link.MaxSize":  406000003,
		"Order.Min":     406000004,
	}
}
