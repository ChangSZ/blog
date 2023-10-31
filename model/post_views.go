package model

type PostViews struct {
	ID     int `gorm:"column:id;primary_key"`
	PostID int `gorm:"column:post_id"`
	Num    int `gorm:"column:num"`
}

func (t *PostViews) TableName() string {
	return "post_views"
}
