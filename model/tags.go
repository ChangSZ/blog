package model

import (
	"time"
)

type Tags struct {
	ID          int       `gorm:"column:id;primary_key"`
	Name        string    `gorm:"column:name"`
	DisplayName string    `gorm:"column:display_name"`
	SeoDesc     string    `gorm:"column:seo_desc"`
	Num         int       `gorm:"column:num"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Tags) TableName() string {
	return "tags"
}
