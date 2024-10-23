package logic

import (
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	sqlba "github.com/bramble555/blog/dao/mysql/banner"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/image"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func UploadImages(c *gin.Context, fileList []*multipart.FileHeader) (*[]image.FileUploadResponse, error) {
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
		size := float64(file.Size) / 1024 / 1024
		if size >= float64(global.Config.Upload.Size) {
			*resFileList = append(*resFileList, image.FileUploadResponse{
				FileName:  file.Filename,
				IsSuccess: false,
				Msg:       fmt.Sprintf("图片太大了,是%.2fMB,图片大小需要缩小到%dMB", size, global.Config.Upload.Size),
			})
			continue
		}
		*resFileList = append(*resFileList, image.FileUploadResponse{
			FileName:  file.Filename,
			IsSuccess: true,
			Msg:       "上传成功",
		})
		err := c.SaveUploadedFile(file, global.Config.Upload.Path+"/"+file.Filename)
		if err != nil {
			global.Log.Errorf("Logic ImageHandler SaveUploadedFile[images] err:%s\n", err.Error())
			continue
		}
		var fileObj multipart.File
		fileObj, err = file.Open()
		if err != nil {
			global.Log.Errorf("Logic file.Open err:%s\n", err.Error())
			continue
		}
		// 记得关闭文件对象
		defer fileObj.Close()
		var byteData []byte
		byteData, err = io.ReadAll(fileObj)
		if err != nil {
			global.Log.Errorf("Logic io.ReadAll err:%s\n", err.Error())
			continue
		}
		ok := sqlba.CheckBannerNotExists(byteData)
		// 如果不存在，插入数据库
		if ok {
			// 写入数据库
			err = sqlba.UploadBanners(byteData, file.Filename)
			if err != nil {
				global.Log.Errorf("sqlba.UploadBanners err:%s\n", err.Error())
			}
		}
	}
	return resFileList, nil
}
func GetBannerList(pl *model.ParamList) ([]model.BannerModel, error) {
	return sqlba.GetBannerList(pl)
}
func DeleteBannerList(pdl *model.ParamDeleteList) (string, error) {
	return sqlba.DeleteBannerList(pdl)
}
