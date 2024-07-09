package console

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/common"
	"github.com/ChangSZ/blog/conf"
	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/infra/strings"
	"github.com/ChangSZ/blog/service"
)

type Post struct {
}

func NewPost() Console {
	return &Post{}
}

func NewPostImg() Img {
	return &Post{}
}

func NewTrash() Trash {
	return &Post{}
}

func (p *Post) Index(ctx *gin.Context) {
	appG := api.Gin{C: ctx}

	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	postList, err := service.ConsolePostIndex(ctx, limit, offset, false)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCount, err := service.ConsolePostCount(ctx, limit, offset, false)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) Create(ctx *gin.Context) {
	cates, err := service.CateListBySort(ctx)
	appG := api.Gin{C: ctx}
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	tags, err := service.AllTags(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = conf.Cnf.ImgUploadUrl
	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) Store(ctx *gin.Context) {
	appG := api.Gin{C: ctx}
	requestJson, exists := ctx.Get("json")
	if !exists {
		log.WithTrace(ctx).Error("get request_params from context fail")
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	var ps common.PostStore
	ps, ok := requestJson.(common.PostStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}

	userId, exist := ctx.Get("userId")
	if !exist || userId.(int) == 0 {
		log.WithTrace(ctx).Error("Can not get user")
		appG.Response(http.StatusOK, 400001004, nil)
		return
	}

	service.PostStore(ctx, ps, userId.(int))
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) Edit(ctx *gin.Context) {
	postIdStr := ctx.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	post, err := service.PostDetail(ctx, postIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postTags, err := service.PostIdTag(ctx, postIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCate, err := service.PostCate(ctx, postIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data := make(map[string]interface{})
	posts := make(map[string]interface{})
	posts["post"] = post
	posts["postCate"] = postCate
	posts["postTag"] = postTags
	data["post"] = posts
	cates, err := service.CateListBySort(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	tags, err := service.AllTags(ctx)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	data["cates"] = cates
	data["tags"] = tags
	data["imgUploadUrl"] = conf.Cnf.ImgUploadUrl
	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) Update(ctx *gin.Context) {
	postIdStr := ctx.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
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
	var ps common.PostStore
	ps, ok := requestJson.(common.PostStore)
	if !ok {
		log.WithTrace(ctx).Error("request_params turn to error")
		appG.Response(http.StatusOK, 400001001, nil)
		return
	}
	service.PostUpdate(ctx, postIdInt, ps)
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) Destroy(ctx *gin.Context) {
	postIdStr := ctx.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	_, err = service.PostDestroy(ctx, postIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) TrashIndex(ctx *gin.Context) {
	appG := api.Gin{C: ctx}

	queryPage := ctx.DefaultQuery("page", "1")
	queryLimit := ctx.DefaultQuery("limit", conf.Cnf.DefaultLimit)

	limit, offset := common.Offset(queryPage, queryLimit)
	postList, err := service.ConsolePostIndex(ctx, limit, offset, true)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	queryPageInt, err := strconv.Atoi(queryPage)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	postCount, err := service.ConsolePostCount(ctx, limit, offset, true)

	data := make(map[string]interface{})
	data["list"] = postList
	data["page"] = common.MyPaginate(postCount, limit, queryPageInt)

	appG.Response(http.StatusOK, 0, data)
	return
}

func (p *Post) UnTrash(ctx *gin.Context) {
	postIdStr := ctx.Param("id")
	postIdInt, err := strconv.Atoi(postIdStr)
	appG := api.Gin{C: ctx}

	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	_, err = service.PostUnTrash(ctx, postIdInt)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}
	appG.Response(http.StatusOK, 0, nil)
	return
}

func (p *Post) ImgUpload(ctx *gin.Context) {
	appG := api.Gin{C: ctx}

	errFiles := []string{}
	succMap := make(map[string]string)

	err := ctx.Request.ParseMultipartForm(32 << 20) // 32 MB limit for the entire request
	if err != nil {
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}

	files := ctx.Request.MultipartForm.File["file[]"]
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			errFiles = append(errFiles, fileHeader.Filename)
			continue
		}
		defer file.Close()

		filename := fileHeader.Filename

		// 为文件生成唯一的名称，可以使用时间戳或其他方法
		saveName := strings.RandString(5) + "_" + filename
		savePath := conf.Cnf.ImgUploadDst + saveName
		out, err := os.Create(savePath)
		if err != nil {
			errFiles = append(errFiles, filename)
			continue
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			errFiles = append(errFiles, filename)
			continue
		}

		if conf.Cnf.QiNiuUploadImg {
			go service.Qiniu(ctx, savePath, saveName)
			succMap[filename] = conf.Cnf.QiNiuHostName + saveName
			continue
		}

		succMap[filename] = conf.Cnf.AppImgUrl + saveName
	}

	data := map[string]interface{}{
		"errFiles": errFiles,
		"succMap":  succMap,
	}

	appG.Response(http.StatusOK, 0, data)
	return
}
