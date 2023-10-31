package console

import (
	"net/http"
	"strconv"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/service"
	"github.com/gin-gonic/gin"
)

type Tag struct {
}

func NewTag() Console {
	return &Tag{}
}

func (t *Tag) Index(ctx *gin.Context) {
	appG := api.Gin{C: ctx}

	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	count, tags, err := service.TagsIndex(ctx, limit, offset)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 402000001, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = tags
	data["page"] = common.MyPaginate(count, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (t *Tag) Create(ctx *gin.Context) {}

func (t *Tag) Store(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	var ts common.TagStore
	ts, ok := requestJson.(common.TagStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err := service.TagStore(ctx, ts)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 403000006, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (t *Tag) Edit(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	tagIdStr := ctx.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}
	tagData, err := service.GetTagById(ctx, tagIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 403000008, nil)
		return
	}
	appG.Response(http.StatusOK, 0, tagData)
	return
}

func (t *Tag) Update(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	tagIdStr := ctx.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
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
	ts, ok := requestJson.(common.TagStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.TagUpdate(ctx, tagIdInt, ts)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 403000007, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (t *Tag) Destroy(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	tagIdStr := ctx.Param("id")
	tagIdInt, err := strconv.Atoi(tagIdStr)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001002, nil)
		return
	}

	_, err = service.GetTagById(ctx, tagIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 403000008, nil)
		return
	}
	service.DelTagRel(ctx, tagIdInt)
	appG.Response(http.StatusOK, 0, nil)
	return
}
