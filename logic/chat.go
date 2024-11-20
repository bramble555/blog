package logic

import (
	"encoding/json"
	"fmt"
	"image/png"
	"os"
	"strings"
	"sync"

	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/gorilla/websocket"
	"github.com/o1egl/govatar"
	"github.com/ser163/WordBot/generate"
)

var (
	ConnGroupMap = make(map[string]*websocket.Conn)
	mu           sync.RWMutex
)

func UploadChat(cm *model.ChatModel) error {
	return global.DB.Create(&cm).Error
}

// AddUserConnection 添加用户连接
func AddUserConnection(addr string, conn *websocket.Conn) {
	mu.Lock()
	defer mu.Unlock()
	ConnGroupMap[addr] = conn
}

// RemoveUserConnection 移除用户连接
func RemoveUserConnection(addr string) {
	mu.Lock()
	defer mu.Unlock()
	delete(ConnGroupMap, addr)
}

// SendUserMsg 发送单个用户消息
func SendUserMsg(addr string, pcg *model.ResponseChatGroup) {
	mu.RLock()
	defer mu.RUnlock()
	byteData, _ := json.Marshal(pcg)
	if conn, exists := ConnGroupMap[addr]; exists {
		conn.WriteMessage(websocket.TextMessage, byteData)
	}
}

// SendGroupMsg 群聊消息
func SendGroupMsg(pcg *model.ResponseChatGroup) {
	mu.RLock()
	defer mu.RUnlock()
	byteData, _ := json.Marshal(pcg)
	for _, conn := range ConnGroupMap {
		conn.WriteMessage(websocket.TextMessage, byteData)
		global.Log.Debugf("消息是%s", string(byteData))
	}
}

// GetIP 提取地址中的 IP
func GetIP(addr string) string {
	return strings.Split(addr, ":")[0]
}

// GetNameAvatar 生成随机名字和头像
func GetNameAvatar() (string, string, error) {
	// 生成随机名字
	w, err := generate.GenRandomWorld(0, "none")
	if err != nil {
		global.Log.Errorf("生成随机世界时出错: %s\n", err.Error())
		return "", "", err
	}

	// 生成头像
	img, err := govatar.GenerateForUsername(govatar.MALE, w.Word)
	if err != nil {
		global.Log.Errorf("生成头像时出错: %s\n", err.Error())
		return "", "", err
	}

	// 保存头像
	filePath := fmt.Sprintf("./uploads/avatar/%s.png", w.Word)
	file, err := os.Create(filePath)
	if err != nil {
		global.Log.Errorf("创建文件时出错: %s\n", err.Error())
		return "", "", err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		global.Log.Errorf("保存图像时出错: %s\n", err.Error())
		return "", "", err
	}
	global.Log.Println("头像生成并保存成功")
	return w.Word, filePath, nil
}
