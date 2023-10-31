package zip

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ChangSZ/blog/infra/file"
	"github.com/ChangSZ/blog/infra/log"
)

// dst: 保存压缩包的本地路径
// dirs: 需要打包的本地文件夹路径, 可以有多个
// CompressDirs 压缩多个文件夹
func CompressDirs(dst string, dirs ...string) error {
	file.MakeDirByFile(dst)
	zipFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer zipFile.Close()
	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, dir := range dirs {
		if !file.FileOrDirExists(dir) {
			log.Warn("文件路径不存在", dir)
			continue
		}
		err := compressDir(dir, archive)
		if err != nil {
			log.Errorf("压缩文件夹时出现错误: %v", err)
			return err
		}
	}
	return nil
}

func compressDir(dir string, archive *zip.Writer) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if path == dir {
			return nil
		}
		info, _ := d.Info()
		h, _ := zip.FileInfoHeader(info)
		h.Name = strings.TrimPrefix(path, dir+"/")
		if info.IsDir() {
			h.Name += "/"
		} else {
			h.Method = zip.Deflate
		}
		writer, _ := archive.CreateHeader(h)
		if !info.IsDir() {
			srcFile, err := os.Open(path)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			io.Copy(writer, srcFile)
		}
		return nil
	})
}
