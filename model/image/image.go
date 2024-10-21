package image

type FileUploadResponse struct {
	FileName  string `json:"file_name"`
	IsSuccess bool   `json:"is_success"`
	Msg       string `json:"msg"`
}

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
