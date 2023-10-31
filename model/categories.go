package model

import (
	"time"
)

type Categories struct {
	ID          int       `gorm:"column:id;primary_key"`
	Name        string    `gorm:"column:name"`
	DisplayName string    `gorm:"column:display_name"`
	SeoDesc     string    `gorm:"column:seo_desc"`
	ParentID    int       `gorm:"column:parent_id"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Categories) TableName() string {
	return "categories"
}
