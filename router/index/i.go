package index

import "github.com/gin-gonic/gin"

type Home interface {
	Index(*gin.Context)
	IndexTag(*gin.Context)
	IndexCate(*gin.Context)
	Detail(*gin.Context)
	Archives(*gin.Context)
	NoFound(*gin.Context)
	Rss(*gin.Context)
	Atom(*gin.Context)
}
