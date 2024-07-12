package model

import (
	"time"
)

type Posts struct {
	ID        int        `gorm:"column:id;primary_key"`
	UID       string     `gorm:"column:uid"`
	UserID    int        `gorm:"column:user_id"`
	Title     string     `gorm:"column:title"`
	Summary   string     `gorm:"column:summary"`
	Content   string     `gorm:"column:content"`
	Password  string     `gorm:"column:password"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time  `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Posts) TableName() string {
	return "posts"
}
