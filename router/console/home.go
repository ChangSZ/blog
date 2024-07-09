package console

import (
	"net/http"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/service"
)

type HomeStatistics struct {
}

func NewStatistics() Statistics {
	return &HomeStatistics{}
}

func (h *HomeStatistics) Index(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	postCnt, err := service.PostCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	cateCnt, err := service.CateCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	tagCnt, err := service.TagCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	linkCnt, err := service.LinkCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	userCnt, err := service.UserCnt(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	var data []interface{}
	pcnt := Res{
		Title: "文章总数",
		Icon:  "ios-book-outline",
		Count: postCnt,
		Color: "#ff9900",
	}
	data = append(data, pcnt)
	ucnt := Res{
		Title: "用户总数",
		Icon:  "md-person-add",
		Count: userCnt,
		Color: "#2d8cf0",
	}
	data = append(data, ucnt)
	lcnt := Res{
		Title: "外链总数",
		Icon:  "ios-link",
		Count: linkCnt,
		Color: "#E46CBB",
	}
	data = append(data, lcnt)
	ccnt := Res{
		Title: "分类总数",
		Icon:  "md-locate",
		Count: cateCnt,
		Color: "#19be6b",
	}
	data = append(data, ccnt)
	tcnt := Res{
		Title: "标签总数",
		Icon:  "md-share",
		Count: tagCnt,
		Color: "#39ed14",
	}
	data = append(data, tcnt)
	qcnt := Res{
		Title: "未知BUG",
		Icon:  "ios-bug",
		Count: 998,
		Color: "#ed3f14",
	}
	data = append(data, qcnt)
	appG.Response(http.StatusOK, 0, data)
	return
}

type Res struct {
	Title string `json:"title"`
	Icon  string `json:"icon"`
	Count int64  `json:"count"`
	Color string `json:"color"`
}
