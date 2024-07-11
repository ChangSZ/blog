package console

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ChangSZ/golib/log"
	"github.com/gin-gonic/gin"

	"github.com/ChangSZ/blog/infra/gin/api"
	"github.com/ChangSZ/blog/service"
)

type Minio struct{}

func NewMinio() MinioI {
	return &Minio{}
}

func (p *Minio) GetFile(ctx *gin.Context) {
	appG := api.Gin{C: ctx}

	type presignedURLRequest struct {
		Bucket     string `form:"bucket"`
		ObjectName string `form:"objectName"`
	}
	req := new(presignedURLRequest)
	if err := ctx.ShouldBind(req); err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 401000004, nil)
		return
	}
	if req.Bucket == "" {
		appG.Response(http.StatusOK, 401000004, "非预期的URL")
		return
	}

	url, err := service.NewMinio().PresignedURL(ctx, req.Bucket, req.ObjectName)
	if err != nil {
		log.WithTrace(ctx).Error(err)
		appG.Response(http.StatusOK, 500000000, nil)
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, fmt.Errorf("failed to fetch file: %v", err))
		return
	}
	defer resp.Body.Close()

	// 检查请求是否成功
	if resp.StatusCode != http.StatusOK {
		appG.Response(http.StatusOK, 500000000, fmt.Errorf("failed to fetch file: %v", err))
		return
	}

	parts := strings.Split(req.ObjectName, "/")
	filename := parts[len(parts)-1]
	// 设置响应头
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", resp.Header.Get("Content-Type"))

	// 将文件内容写入响应
	_, err = io.Copy(ctx.Writer, resp.Body)
	if err != nil {
		appG.Response(http.StatusOK, 500000000, fmt.Errorf("failed to write file: %v", err))
		return
	}
	appG.Response(http.StatusOK, 0, nil)
}
