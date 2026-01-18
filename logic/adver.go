package logic

import (
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/bramble555/blog/dao/mysql/advert"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

func CreateAdvert(ad *model.AdvertModel) (string, error) {
	return advert.CreateAdvert(ad)
}

func GetAdvertList(pl *model.ParamList, isShow bool) (*model.PageResult[model.AdvertModel], error) {
	list, count, err := advert.GetAdvertList(pl, isShow)
	if err != nil {
		return nil, err
	}
	return &model.PageResult[model.AdvertModel]{
		List:  list,
		Count: count,
	}, nil
}

func DeleteAdvertList(pdl *model.ParamDeleteList) (string, error) {
	return advert.DeleteAdvertList(pdl)
}

func UpdateAdvertShow(p *model.ParamUpdateAdvertShow) (string, error) {
	return advert.UpdateAdvertShow(p)
}

// UploadAdvertImages 处理广告图片上传，文件保存到配置的 AdPath 目录，并返回可直接使用的 URL 路径
func UploadAdvertImages(c *gin.Context, fileList []*multipart.FileHeader) (*[]model.FileUploadResponse, error) {
	resFileList := new([]model.FileUploadResponse)
	if len(fileList) == 0 {
		return resFileList, nil
	}

	pkg.CreateFolder(global.Config.Upload.AdPath)

	for _, file := range fileList {
		fileExt := strings.Split(file.Filename, ".")
		if len(fileExt) != 2 {
			global.Log.Errorf("上传的文件%s没有扩展名\n", file.Filename)
			continue
		}
		if _, exists := model.WhiteImageExtList[fileExt[1]]; !exists {
			global.Log.Errorf("上传的文件%s的扩展名不被支持,文件名是%s\n", fileExt[1], file.Filename)
			continue
		}

		size := float64(file.Size) / 1024 / 1024
		if size >= float64(global.Config.Upload.Size) {
			*resFileList = append(*resFileList, model.FileUploadResponse{
				FileName: "",
				Msg:      fmt.Sprintf("图片太大了,是%.2fMB,图片大小需要缩小到%dMB\n", size, global.Config.Upload.Size),
			})
			continue
		}

		targetPath := global.Config.Upload.AdPath + "/" + file.Filename
		if err := c.SaveUploadedFile(file, targetPath); err != nil {
			global.Log.Errorf("UploadAdvertImages SaveUploadedFile err:%s\n", err.Error())
			continue
		}

		urlPath := "/uploads/ad/" + file.Filename
		*resFileList = append(*resFileList, model.FileUploadResponse{
			FileName: urlPath,
		})
	}

	return resFileList, nil
}
