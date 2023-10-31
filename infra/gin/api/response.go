package api

import (
	"github.com/ChangSZ/blog/conf"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type ds struct{}

func (g *Gin) Response(httpCode, errCode int, data interface{}) {
	if data == nil {
		data = ds{}
	}
	msg := conf.GetMsg(errCode)
	g.C.JSON(httpCode, gin.H{
		"code":    errCode,
		"message": msg,
		"data":    data,
	})
	g.C.Abort()
	return
}
