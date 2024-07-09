package index

import (
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/model"
	"github.com/ChangSZ/blog/service"
)

type Web struct {
	ApiController
}

func NewIndex() Home {
	return &Web{}
}

func (w *Web) Index(ctx *gin.Context) {
	w.C = ctx
	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)

	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(ctx, queryPage, queryLimit, "default", "")
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["title"] = h["system"].(*model.Systems).Title
	w.Response(http.StatusOK, 0, h)
	return
}

func (w *Web) IndexTag(ctx *gin.Context) {
	w.C = ctx
	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)
	name := ctx.Param("name")
	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(ctx, queryPage, queryLimit, "tag", name)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["tagName"] = name
	h["tem"] = "tagList"
	h["title"] = template.HTML(name + " --  tags &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*model.Systems).Title)

	ctx.HTML(http.StatusOK, "master.go.tmpl", h)
	return
}

func (w *Web) IndexCate(ctx *gin.Context) {
	w.C = ctx
	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultIndexLimit)
	name := ctx.Param("name")

	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	postData, err := service.IndexPost(ctx, queryPage, queryLimit, "cate", name)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["post"] = postData.PostListArr
	h["paginate"] = postData.Paginate
	h["cateName"] = name
	h["tem"] = "cateList"
	h["title"] = template.HTML(name + " --  category &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*model.Systems).Title)

	w.Response(http.StatusOK, 0, h)
	return
}

func (w *Web) Detail(ctx *gin.Context) {
	w.C = ctx
	postIdStr := ctx.Param("id")

	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	postDetail, err := service.IndexPostDetail(ctx, postIdStr)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	go service.PostViewAdd(ctx, postIdStr)

	github := common.IndexGithubParam{
		GithubName:         conf.Cnf.GithubName,
		GithubRepo:         conf.Cnf.GithubRepo,
		GithubClientId:     conf.Cnf.GithubClientId,
		GithubClientSecret: conf.Cnf.GithubClientSecret,
		GithubLabels:       conf.Cnf.GithubLabels,
	}

	h["post"] = postDetail
	h["github"] = github
	h["tem"] = "detail"
	h["title"] = template.HTML(postDetail.Post.Title + " &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*model.Systems).Title)

	w.Response(http.StatusOK, 0, h)
	return
}

func (w *Web) Archives(ctx *gin.Context) {
	w.C = ctx
	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	res, err := service.PostArchives(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")

	var dateIndexs []int
	for k := range res {
		tt, _ := time.ParseInLocation("2006-01-02 15:04:05", k+"-01 00:00:00", loc)
		dateIndexs = append(dateIndexs, int(tt.Unix()))
	}

	sort.Sort(sort.Reverse(sort.IntSlice(dateIndexs)))

	var newData []interface{}
	for _, j := range dateIndexs {
		dds := make(map[string]interface{})
		tm := time.Unix(int64(j), 0)
		dateIndex := tm.Format("2006-01")
		dds["dates"] = dateIndex
		dds["lists"] = res[dateIndex]
		newData = append(newData, dds)
	}

	h["tem"] = "archives"
	h["archives"] = newData
	h["title"] = template.HTML("归档 &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*model.Systems).Title)

	w.Response(http.StatusOK, 0, h)
	return
}

func (w *Web) Rss(ctx *gin.Context) {
	w.C = ctx
	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	feed, err := service.CommonRss(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	rss, err := feed.ToRss()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}
	h["rss"] = rss
	h["title"] = template.HTML("Rss &nbsp;&nbsp;-&nbsp;&nbsp;" + h["system"].(*model.Systems).Title)
	w.Response(http.StatusOK, 528, h)
	return
}

func (w *Web) Atom(ctx *gin.Context) {
	w.C = ctx
	h, err := service.CommonData(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	feed, err := service.CommonRss(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	atom, err := feed.ToAtom()
	if err != nil {
		log.WithTrace(ctx).Error(err)
		w.Response(http.StatusOK, 404, h)
		return
	}

	h["atom"] = atom
	w.Response(http.StatusOK, 633, h)
	return
}

func (w *Web) NoFound(ctx *gin.Context) {
	w.C = ctx
	w.Response(http.StatusOK, 404, gin.H{
		"themeJs":  "/static/assets/js",
		"themeCss": "/static/assets/css",
	})
	return
}
