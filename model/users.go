package model

import (
	"time"
)

type Users struct {
	ID              int       `gorm:"column:id;primary_key"`
	Name            string    `gorm:"column:name"`
	Email           string    `gorm:"column:email"`
	Status          int       `gorm:"column:status"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at"`
	Password        string    `gorm:"column:password"`
	RememberToken   string    `gorm:"column:remember_token"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (t *Users) TableName() string {
	return "users"
}
