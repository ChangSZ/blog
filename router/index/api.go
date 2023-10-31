package index

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	C *gin.Context
}

func (a *ApiController) Response(httpCode, errCode int, data gin.H) {
	if data == nil {
		panic("常规信息应该设置")
	}
	if errCode == 500 {
		a.C.HTML(http.StatusOK, "5xx.go.tmpl", data)
	} else if errCode == 404 {
		a.C.HTML(http.StatusOK, "4xx.go.tmpl", data)
	} else if errCode == 0 {
		a.C.HTML(http.StatusOK, "master.go.tmpl", data)
	} else if errCode == 528 {
		a.C.HTML(http.StatusOK, "rss.go.tmpl", data)
	} else if errCode == 633 {
		a.C.HTML(http.StatusOK, "atom.go.tmpl", data)
	} else {
		a.C.HTML(http.StatusOK, "5xx.go.tmpl", nil)
	}

	a.C.Abort()
	return
}
