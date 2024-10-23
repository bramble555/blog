package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

func validateListParams(c *gin.Context) (*model.ParamList, error) {
	// 默认值
	pl := new(model.ParamList)
	pl.Page = 1
	pl.Size = 10
	pl.Order = model.OrderByTime

	// 绑定查询参数
	err := c.ShouldBindQuery(pl)
	if err != nil {
		global.Log.Errorf("controller BindAndValidateParams ShouldBindQuery err:%s\n", err.Error())
		ResponseErrorWithData(c, CodeInvalidParam, err.Error())
		return pl, err
	}

	// 参数校验
	if pl.Page < 1 {
		pl.Page = 1
	}
	if pl.Size <= 0 {
		pl.Size = 10
	}
	if pl.Order == "" {
		pl.Order = model.OrderByTime
	}
	return pl, nil
}
