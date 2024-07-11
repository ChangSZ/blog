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

type Category struct{}

func NewCategory() Console {
	return &Category{}
}

func (cate *Category) Index(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	cates, err := service.CateListBySort(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000001, nil)
		return
	}
	appG.Response(http.StatusOK, 0, cates)
}

func (cate *Category) Create(ctx *gin.Context) {}

func (cate *Category) Store(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var cs common.CateStore
	cs, ok := requestJson.(common.CateStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	_, err := service.CateStore(ctx, cs)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000010, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
}

func (cate *Category) Edit(ctx *gin.Context) {
	cateIdStr := ctx.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	cateData, err := service.GetCateById(ctx, cateIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, cateData)
}

func (cate *Category) Update(ctx *gin.Context) {
	cateIdStr := ctx.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var cs common.CateStore
	cs, ok := requestJson.(common.CateStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	_, err = service.CateUpdate(ctx, cateIdInt, cs)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000009, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
}

func (cate *Category) Destroy(ctx *gin.Context) {
	cateIdStr := ctx.Param("id")
	cateIdInt, err := strconv.Atoi(cateIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}

	_, err = service.GetCateById(ctx, cateIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000000, nil)
		return
	}

	pd, err := service.GetCateByParentId(ctx, cateIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000000, nil)
		return
	}
	if pd.ID > 0 {
		log.WithTrace(ctx).Error(err, ", It has children node")
		appG.Response(http.StatusOK, 402000011, nil)
		return
	}

	service.DelCateRel(ctx, cateIdInt)
	appG.Response(http.StatusOK, 0, nil)
}
