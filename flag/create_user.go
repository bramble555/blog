package flag

import (
	sys_flag "flag"
	"fmt"
	"os"

	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model/ctype"
)

func FlagUserParse() {
	// 用户权限，默认是 用户权限
	permission := sys_flag.String("u", "user", "创建用户权限")
	username := sys_flag.String("n", "", "用户名")
	password := sys_flag.String("p", "", "密码")
	sys_flag.Parse()
	if *username == "" {
		fmt.Println("请输入用户名")
		os.Exit(1)
	}
	if *password == "" {
		fmt.Println("请输入密码")
		os.Exit(1)
	}
	// 默认是用户
	role := uint(ctype.PermissionUser)
	// 如果是超级用户，那就 改为超级用户
	if *permission == "admin" {
		role = uint(ctype.PermissionAdmin)
	}
	if *permission != "admin" && *permission != "user" {
		fmt.Println("-u 输入参数有误")
		os.Exit(1)
	}
	// 用命令行创建用户
	err := user.CreateFlagUser(role, *username, *password)
	if err != nil {
		global.Log.Errorf("user CreateFlagUser err: %s\n", err.Error())
		os.Exit(1)
	}
	global.Log.Infof("创建用户成功\n")
	os.Exit(0)
}
