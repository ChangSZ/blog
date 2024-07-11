package console

import (
	"net/http"
	"strconv"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/service"
)

type Home struct {
}

func NewHome() System {
	return &Home{}
}

func (s *Home) Index(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	themes := make(map[int]interface{})
	themes[1] = 1
	system, err := service.GetSystemList(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		return
	}
	data := make(map[string]interface{})
	data["themes"] = themes
	data["system"] = system
	log.WithTrace(ctx).Info(" Succeed to get system index ")
	appG.Response(http.StatusOK, 0, data)
}

func (s *Home) Update(ctx *gin.Context) {
	systemIdStr := ctx.Param("id")
	systemIdInt, err := strconv.Atoi(systemIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
	}

	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	//var ss common.ConsoleSystem
	ss, ok := requestJson.(common.ConsoleSystem)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.SystemUpdate(ctx, systemIdInt, ss)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 405000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
}
