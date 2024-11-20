package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/logic"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func GetChatGroupHandler(c *gin.Context) {
	// 升级为 websocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// true 表示放行
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		global.Log.Errorf("upgrade websocket err:%s\n", err.Error())
		ResponseError(c, CodeServerBusy)
		return
	}
	addr := conn.RemoteAddr().String()
	global.Log.Debugf("%s连接成功", addr)

	// 添加到连接组
	logic.AddUserConnection(addr, conn)
	defer func() {
		// 移除连接
		logic.RemoveUserConnection(addr)
		conn.Close()
	}()

	// 生成随机昵称和头像
	name, avatar, err := logic.GetNameAvatar()
	if err != nil {
		ResponseError(c, CodeServerBusy)
		return
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			// 用户断开
			break
		}

		// 参数绑定
		var pcg model.ParamChatGroup
		if err := json.Unmarshal(p, &pcg); err != nil {
			continue
		}

		// 设置默认昵称和头像
		if pcg.NickName == "" {
			pcg.NickName = name
		}
		if pcg.Avatar == "" {
			pcg.Avatar = avatar
		}

		// 初始化响应结构
		response := model.ResponseChatGroup{
			ParamChatGroup: pcg,
			Date:           time.Now(),
		}
		global.Log.Debugf("response:%+v", response)
		// 根据消息类型处理
		if pcg.MsgType == ctype.TextMsg || pcg.MsgType == ctype.ImageMsg {
			// 处理测试消息或图片消息
			logic.SendGroupMsg(&response)
		} else if pcg.MsgType == ctype.InRoomMsg {
			// 处理进入消息
			response.ParamChatGroup.Content = fmt.Sprintf("欢迎%s来到聊天室", pcg.NickName)
			logic.SendGroupMsg(&response)
		} else {
			// 处理未知消息类型
			pcg.Content = "消息类型错误"
			response.ParamChatGroup = pcg
			logic.SendUserMsg(addr, &response)
		}

		// 保存聊天记录
		chatModel := model.ChatModel{
			NickName: response.ParamChatGroup.NickName,
			Avatar:   response.ParamChatGroup.Avatar,
			Content:  response.ParamChatGroup.Content,
			IP:       logic.GetIP(addr),
			Addr:     addr,
			MsgType:  response.ParamChatGroup.MsgType,
		}
		if err := logic.UploadChat(&chatModel); err != nil {
			global.Log.Errorf("上传聊天记录失败: %s", err.Error())
			ResponseError(c, CodeServerBusy)
			return
		}
	}
}
