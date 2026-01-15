package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/bramble555/blog/dao/mysql/banner"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/pkg"
	"github.com/gin-gonic/gin"
)

// 上传多个文件
func UploadImages(c *gin.Context, fileList []*multipart.FileHeader) (*[]model.FileUploadResponse, error) {
	resFileList := new([]model.FileUploadResponse)
	pkg.CreateFolder(global.Config.Upload.Path)
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
				FileName: file.Filename,
				Msg:      fmt.Sprintf("图片太大了,是%.2fMB,图片大小需要缩小到%dMB\n", size, global.Config.Upload.Size),
			})
			continue
		}
		fileObj, err := file.Open()
		if err != nil {
			global.Log.Errorf("Logic file.Open err:%s\n", err.Error())
			continue
		}
		// 需要关闭文件流
		defer fileObj.Close()
		hash := md5.New()
		// 使用 io.Copy 将文件内容拷贝到 hash 中,而不是 io.ReadAll,
		// 因为 io.ReadAll 会将文件内容全部读取到内存中,如果文件很大,会导致内存溢出
		if _, err := io.Copy(hash, fileObj); err != nil {
			global.Log.Errorf("Logic io.Copy err:%s\n", err.Error())
			continue
		}
		byteHash := hash.Sum(nil)
		hashStr := hex.EncodeToString(byteHash)
		ok := banner.CheckBannerNotExists(hashStr)
		if !ok {
			*resFileList = append(*resFileList, model.FileUploadResponse{
				FileName: file.Filename,
				Msg:      "该图片已上传",
			})
			continue
		}
		// 与前端商量好,如果上传成功,msg 就为 "",否则 msg 就有内容
		*resFileList = append(*resFileList, model.FileUploadResponse{
			FileName: file.Filename,
		})
		err = c.SaveUploadedFile(file, global.Config.Upload.Path+"/"+file.Filename)
		if err != nil {
			global.Log.Errorf("Logic BannerHandler SaveUploadedFile[banners] err:%s\n", err.Error())
			continue
		}
		err = banner.UploadBanners(hashStr, file.Filename)
		if err != nil {
			global.Log.Errorf("sqlba.UploadBanners err:%s\n", err.Error())
		}
	}
	return resFileList, nil
}
func GetBannerList(pl *model.ParamList) (*[]model.BannerModel, error) {
	return banner.GetBannerList(pl)
}

func DeleteBannerList(pdl *model.ParamDeleteList) (string, error) {
	return banner.DeleteBannerList(pdl)
}
