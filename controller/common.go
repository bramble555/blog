package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/gin-gonic/gin"
)

// validateListParams 验证列表参数
func validateListParams(c *gin.Context) (*model.ParamList, error) {
	pl := new(model.ParamList)
	// 绑定查询参数
	err := c.ShouldBindQuery(pl)
	if err != nil {
		global.Log.Errorf("controller validateListParams ShouldBindQuery err:%s\n", err.Error())
		return pl, err
	}

	// 参数校验 & 处理默认值
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
