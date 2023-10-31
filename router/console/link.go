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

type Link struct{}

func NewLink() Console {
	return &Link{}
}

func (l *Link) Index(ctx *gin.Context) {
	appG := api.Gin{C: ctx}

	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultLimit)
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	limit, offset := common.Offset(queryPage, queryLimit)

	links, cnt, err := service.LinkList(ctx, offset, limit)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["list"] = links
	data["page"] = common.MyPaginate(cnt, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (l *Link) Create(ctx *gin.Context) {}

func (l *Link) Store(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	ls, ok := requestJson.(common.LinkStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	err := service.LinkSore(ctx, ls)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 406000005, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (l *Link) Edit(ctx *gin.Context) {
	linkIdStr := ctx.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	link, err := service.LinkDetail(ctx, linkIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 406000006, nil)
		return
	}
	appG.Response(http.StatusOK, 0, link)
	return
}

func (l *Link) Update(ctx *gin.Context) {
	linkIdStr := ctx.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 400001003, nil)
		return
	}
	ls, ok := requestJson.(common.LinkStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	err = service.LinkUpdate(ctx, ls, linkIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 406000007, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (l *Link) Destroy(ctx *gin.Context) {
	linkIdStr := ctx.Param("id")
	linkIdInt, err := strconv.Atoi(linkIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	err = service.LinkDestroy(ctx, linkIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}
