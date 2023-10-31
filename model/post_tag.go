package model

type PostTag struct {
	ID     int `gorm:"column:id;primary_key"`
	PostID int `gorm:"column:post_id"`
	TagID  int `gorm:"column:tag_id"`
}

func (t *PostTag) TableName() string {
	return "post_tag"
}
