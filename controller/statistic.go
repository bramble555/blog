package controller

import (
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func GetUserLoginHandler(c *gin.Context) {
	if _claims, ok := c.Get("claims"); ok {
		claims := _claims.(*jwt.MyClaims)
		global.Log.Infof("login stats request user_sn=%d role=%d", claims.SN, claims.Role)
	} else {
		global.Log.Warnf("login stats request claims missing")
	}
	data, err := logic.GetUserLoginData()
	if err != nil {
		global.Log.Errorf("controller GetUserLoginHandler logic.GetUserLoginData err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	global.Log.Infof("login stats response items=%d", len(data))
	ResponseSucceed(c, data)
}
func GetDataSumHandler(c *gin.Context) {
	data, err := logic.GetDataSum()
	if err != nil {
		global.Log.Errorf("controller GetDataSumHandler logic.GetDataSum err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	ResponseSucceed(c, data)
}
