package model

type PostCate struct {
	ID     int `gorm:"column:id;primary_key"`
	PostID int `gorm:"column:post_id"`
	CateID int `gorm:"column:cate_id"`
}

func (t *PostCate) TableName() string {
	return "post_cate"
}
