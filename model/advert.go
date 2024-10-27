package model

type AdvertModel struct {
	*MODEL
	Title  string `json:"title" binding:"required"`  // 显示的标题
	Href   string `json:"href" binding:"required"`   // 跳转链接
	Images string `json:"images" binding:"required"` // 图片
	IsShow bool   `json:"is_show"`                   // 是否展示 不能写 required 因为 false 是默认值
}
