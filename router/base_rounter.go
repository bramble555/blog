package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitBaseRoutes(r *gin.RouterGroup) gin.IRoutes {
	base := r.Group("/base")
	fmt.Println(base)
	base.GET("/ping")
	// {
	// 	base.GET("ping", controller.Demo)
	// 	base.GET("encryptpwd", controller.Base.EncryptPasswd) // 生成加密密码
	// 	base.GET("decryptpwd", controller.Base.DecryptPasswd) // 密码解密为明文
	// 	// 登录登出刷新token无需鉴权
	// 	base.POST("/login", authMiddleware.LoginHandler)
	// 	base.POST("/logout", authMiddleware.LogoutHandler)
	// 	base.POST("/refreshToken", authMiddleware.RefreshHandler)
	// 	base.POST("/sendcode", controller.Base.SendCode)   // 给用户邮箱发送验证码
	// 	base.POST("/changePwd", controller.Base.ChangePwd) // 修改用户密码
	// 	base.GET("/dashboard", controller.Base.Dashboard)  // 系统首页展示数据
	// }
	return r
}
