package model

type TagModel struct {
	MODEL
	Title string `json:"title" binding:"required"` //标签名称
}

func (TagModel) TableName() string {
	return "tag_models"
}
