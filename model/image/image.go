package image

type FileUploadResponse struct {
	FileName  string `json:"file_name"`
	IsSuccess bool   `json:"is_success"`
	Msg       string `json:"msg"`
}
type ParamImageList struct {
	Page  int    `json:"page" form:"page"`
	Size  int    `json:"size" form:"size"`
	Order string `json:"order" form:"order"`
}


// 默认按照创建时间降序排序
const OrderByTime = "create_time DESC"

var WhiteImageExtList = map[string]struct{}{
	"jpg":  {},
	"jpeg": {},
	"png":  {},
	"ico":  {},
	"tiff": {},
	"gif":  {},
	"svg":  {},
	"webg": {},
}
