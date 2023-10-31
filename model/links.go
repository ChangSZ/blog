package model

import (
	"time"
)

type Links struct {
	ID        int       `gorm:"column:id;primary_key"`
	Name      string    `gorm:"column:name"`
	Link      string    `gorm:"column:link"`
	Order     int       `gorm:"column:order"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Links) TableName() string {
	return "links"
}
