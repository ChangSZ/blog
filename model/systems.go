package model

import (
	"time"
)

type Systems struct {
	ID           int       `gorm:"column:id;primary_key"`
	Theme        int       `gorm:"column:theme"`
	Title        string    `gorm:"column:title"`
	Keywords     string    `gorm:"column:keywords"`
	Description  string    `gorm:"column:description"`
	RecordNumber string    `gorm:"column:record_number"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Systems) TableName() string {
	return "systems"
}
