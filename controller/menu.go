package controller

import (
	"strconv"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

// UploadMenuHandler 一次只能上传一个菜单
func UploadMenuHandler(c *gin.Context) {
	mm := model.MenuModel{}
	mm.BannerID = nil
	err := c.ShouldBindJSON(&mm)
	if err != nil {
		global.Log.Errorf("controller UploadMenuHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	data, err := logic.UploadMenu(&mm)
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	ResponseSucceed(c, data)

}
func GetMenuListHandler(c *gin.Context) {
	data, err := logic.GetMenuList()
	if err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	ResponseSucceed(c, data)
}
func UpdateMenuHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		global.Log.Errorf("id 有误 err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidID, codeMsgMap[CodeInvalidID])
		return
	}
	mm := model.MenuModel{}
	err = c.ShouldBindJSON(&mm)
	if err != nil {
		global.Log.Errorf("controller UpdateMenuHandler ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.UpdateMenu(uint(id), &mm)
	if err != nil {
		global.Log.Errorf("logic UpdateMenu err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)

}
func DeleteMenuListHander(c *gin.Context) {
	pdl := model.ParamDeleteList{}
	err := c.ShouldBindJSON(&pdl)
	if err != nil {
		global.Log.Errorf("controller DeleteMenuListHander ShouldBindJSON err:%s\n", err.Error())
		ResponseError(c, CodeInvalidParam)
		return
	}
	var data string
	data, err = logic.DeleteMenuList(&pdl)
	if err != nil {
		return
	}
	ResponseSucceed(c, data)
}
