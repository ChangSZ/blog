package router

import (
	"html/template"
	"net/http"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/infra/log"
	"github.com/ChangSZ/blog/middleware"
	"github.com/ChangSZ/blog/router/auth"
	"github.com/ChangSZ/blog/router/console"
	"github.com/ChangSZ/blog/router/index"
	"github.com/ChangSZ/blog/validate"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
)

func RoutersInit() *gin.Engine {
	if conf.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()

	r.Use(kgin.Middlewares(recovery.Recovery(), tracing.Server(), logging.Server(log.GetLoggerWithTrace()), middleware.AddTraceCtx))
	r.Use(middleware.CheckExist())

	r.Static("/static/uploads/images/", "./static/uploads/images/")
	consolePost := console.NewPost()
	consoleCate := console.NewCategory()
	consoleTag := console.NewTag()
	postImg := console.NewPostImg()
	trash := console.NewTrash()
	consoleSystem := console.NewHome()
	consoleLink := console.NewLink()
	consoleAuth := auth.NewAuth()
	consoleHome := console.NewStatistics()
	c := r.Group("/console")
	{
		p := c.Group("/post")
		{
			postV := validate.NewValidate().NewPostV.MyValidate()
			p.GET("/", middleware.Permission("console.post.index"), consolePost.Index)
			p.GET("/create", middleware.Permission("console.post.create"), consolePost.Create)
			p.POST("/", middleware.Permission("console.post.store"), postV, consolePost.Store)
			p.GET("/edit/:id", middleware.Permission("console.post.edit"), consolePost.Edit)
			p.PUT("/:id", middleware.Permission("console.post.update"), postV, consolePost.Update)
			p.DELETE("/:id", middleware.Permission("console.post.destroy"), consolePost.Destroy)
			p.GET("/trash", middleware.Permission("console.post.trash"), trash.TrashIndex)
			p.PUT("/:id/trash", middleware.Permission("console.post.unTrash"), trash.UnTrash)

			p.POST("/imgUpload", middleware.Permission("console.post.imgUpload"), postImg.ImgUpload)
		}
		cate := c.Group("/cate")
		{
			cateV := validate.NewValidate().NewCateV.MyValidate()
			cate.GET("/", middleware.Permission("console.cate.index"), consoleCate.Index)
			cate.GET("/edit/:id", middleware.Permission("console.cate.edit"), consoleCate.Edit)
			cate.PUT("/:id", middleware.Permission("console.cate.update"), cateV, consoleCate.Update)
			cate.POST("/", middleware.Permission("console.cate.store"), cateV, consoleCate.Store)
			cate.DELETE("/:id", middleware.Permission("console.cate.destroy"), consoleCate.Destroy)
		}
		tag := c.Group("/tag")
		{
			tagV := validate.NewValidate().NewTagV.MyValidate()
			tag.GET("/", middleware.Permission("console.tag.index"), consoleTag.Index)
			tag.POST("/", middleware.Permission("console.tag.store"), tagV, consoleTag.Store)
			tag.GET("/edit/:id", middleware.Permission("console.tag.edit"), consoleTag.Edit)
			tag.PUT("/:id", middleware.Permission("console.tag.update"), tagV, consoleTag.Update)
			tag.DELETE("/:id", middleware.Permission("console.tag.destroy"), consoleTag.Destroy)
		}
		system := c.Group("/system")
		{
			systemV := validate.NewValidate().NewSystemV.MyValidate()
			system.GET("/", middleware.Permission("console.system.index"), consoleSystem.Index)
			system.PUT("/:id", middleware.Permission("console.system.update"), systemV, consoleSystem.Update)
		}
		link := c.Group("/link")
		{
			linkV := validate.NewValidate().NewLinkV.MyValidate()
			link.GET("/", middleware.Permission("console.link.index"), consoleLink.Index)
			link.POST("/", middleware.Permission("console.link.store"), linkV, consoleLink.Store)
			link.GET("/edit/:id", middleware.Permission("console.link.edit"), consoleLink.Edit)
			link.PUT("/:id", middleware.Permission("console.link.update"), linkV, consoleLink.Update)
			link.DELETE("/:id", middleware.Permission("console.link.destroy"), consoleLink.Destroy)
		}
		c.DELETE("/logout", middleware.Permission("console.auth.logout"), consoleAuth.Logout)
		c.DELETE("/cache", middleware.Permission("console.auth.cache"), consoleAuth.DelCache)
		h := c.Group("/home")
		{
			h.GET("/", middleware.Permission("console.home.index"), consoleHome.Index)
		}

		// 不需要登录状态权限

		al := c.Group("/login")
		{
			authLoginV := validate.NewValidate().NewAuthLoginV.MyValidate()
			al.GET("/", consoleAuth.Login)
			al.POST("/", authLoginV, consoleAuth.AuthLogin)
		}
		ar := c.Group("/register")
		{
			authRegisterV := validate.NewValidate().NewAuthRegister.MyValidate()
			ar.GET("/", consoleAuth.Register)
			ar.POST("/", authRegisterV, consoleAuth.AuthRegister)
		}
	}

	web := index.NewIndex()
	h := r.Group("")
	{
		r.SetFuncMap(template.FuncMap{
			"rem":    common.Rem,
			"MDate":  common.MDate,
			"MDate2": common.MDate2,
		})
		r.LoadHTMLGlob("template/*.go.tmpl")

		r.Static("/static/assets/", "./static/assets/")
		h.GET("/", web.Index)
		h.GET("/categories/:name", web.IndexCate)
		h.GET("/tags/:name", web.IndexTag)
		h.GET("/detail/:id", web.Detail)
		h.GET("/archives", web.Archives)
		h.GET("/rss", web.Rss)
		h.GET("/atom", web.Atom)
		h.GET("/404", web.NoFound)
		h.GET("/favicon.ico", func(ctx *gin.Context) {
			ctx.File("./static/assets/img/favicon.ico")
		})
	}
	return r
}

func recoverHandler(ctx *gin.Context, err interface{}) {
	apiG := api.Gin{C: ctx}
	apiG.Response(http.StatusOK, 400000000, []string{})
	return
}
