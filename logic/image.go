package logic

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model/image"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func UploadImage(c *gin.Context, fileList []*multipart.FileHeader) (*[]image.FileUploadResponse, error) {
	resFileList := new([]image.FileUploadResponse)
	pkg.CreateFolder(global.Config.Upload.Path)
	for _, file := range fileList {
		// 检验扩展名
		fileExt := strings.Split(file.Filename, ".")
		if len(fileExt) != 2 {
			global.Log.Errorf("上传的文件%s没有扩展名", file.Filename)
			continue
		}
		// 检查扩展名是否在白名单中
		if _, exists := image.WhiteImageExtList[fileExt[1]]; !exists {
			global.Log.Errorf("上传的文件%s的扩展名不被支持,文件名是%s", fileExt[1], file.Filename)
			continue
		}
		global.Log.Debugf("fileName:%s\n", file.Filename)
		err := c.SaveUploadedFile(file, global.Config.Upload.Path+"/"+file.Filename)
		if err != nil {
			global.Log.Errorf("Controller ImageHandler SaveUploadedFile[images] err:%s\n", err.Error())
			continue
		}
		size := float64(file.Size) / 1024 / 1024
		if size >= float64(global.Config.Upload.Size) {
			*resFileList = append(*resFileList, image.FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       fmt.Sprintf("图片太大了,是%.2fMB,图片大小需要缩小到%dMB", size, global.Config.Upload.Size),
			})
		} else {
			*resFileList = append(*resFileList, image.FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: true,
				Msg:       "上传成功",
			})
		}
	}
	return resFileList, nil
}
