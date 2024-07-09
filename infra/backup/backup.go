package backup

import (
	"time"

	"github.com/ChangSZ/golib/log"
	"github.com/robfig/cron"

	"github.com/ChangSZ/golib/zip"

	"github.com/ChangSZ/golib/mail"

	"github.com/ChangSZ/blog/infra/conf"
	"github.com/ChangSZ/blog/infra/conn"
)

func (bp *BackUpParam) SetFiles(files ...string) *BackUpParam {
	bp.Files = files
	return bp
}

func (bp *BackUpParam) SetDest(d string) *BackUpParam {
	bp.Dest = d
	return bp
}

func (bp *BackUpParam) SetCronSpec(d string) *BackUpParam {
	bp.CronSpec = d
	return bp
}

func (bp *BackUpParam) SetFileName(fn string) *BackUpParam {
	bp.FileName = fn
	return bp
}

func (bp *BackUpParam) SetFilePath(fp string) *BackUpParam {
	bp.FilePath = fp
	return bp
}

type BackUpParam struct {
	Files    []string         `json:"files"`
	CronSpec string           `json:"cronSpec"`
	Dest     string           `json:"dest"`
	FileName string           `json:"file_name"`
	FilePath string           `json:"file_path"`
	Ep       *mail.EmailParam `json:"ep"`
}

func (bp *BackUpParam) FilePathIsNull() *BackUpParam {
	if bp.FilePath == "" {
		log.Warn("data is null")
		bp.SetFilePath(conf.BackUpFilePath)
	}
	return bp
}

func (bp *BackUpParam) DestIsNull() *BackUpParam {
	if bp.Dest == "" {
		log.Warn("data is null")
		bp.SetDest(conf.BackUpDest)
	}
	return bp
}

func (bp *BackUpParam) FileNameIsNull() *BackUpParam {
	if bp.FileName == "" {
		log.Warn("data is null")
		bp.SetFileName(time.Now().Format("2006-01-02") + conf.BackUpSqlFileName)
	}
	return bp
}

func (bp *BackUpParam) DurationIsNull() *BackUpParam {
	if bp.CronSpec == "" {
		log.Warn("data is null")
		bp.SetCronSpec(conf.BackUpDuration)
	}
	return bp
}

func (bp *BackUpParam) Backup() {
	bp.DestIsNull().FileNameIsNull().FilePathIsNull().DurationIsNull()

	c := cron.New()
	_ = c.AddFunc(bp.CronSpec, bp.doBackUp)
	c.Start()
}

func (bp *BackUpParam) doBackUp() {
	err := conn.MySQLDump(bp.FileName, bp.FilePath)
	if err != nil {
		log.Errorf("back up sql dump is error: %v", err)
	}

	err = zip.CompressDirs(bp.Dest, bp.Files...)
	if err != nil {
		log.Errorf("back up compress is error: %v", err)
		return
	}
	err = bp.Ep.Send()
	if err != nil {
		log.Errorf("back up send mail is error: %v", err)
		return
	}
}
