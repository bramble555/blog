package logic

import (
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bramble555/blog/dao/mysql/chat"
	"github.com/bramble555/blog/dao/mysql/user"
	"github.com/bramble555/blog/global"
	"github.com/bramble555/blog/model"
	"github.com/bramble555/blog/model/ctype"
	"github.com/bramble555/blog/pkg/file"
	"github.com/bramble555/blog/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type chatGroupClient struct {
	conn      *websocket.Conn
	send      chan []byte
	room      *chatGroupRoom
	userSN    int64
	nickName  string
	avatar    string
	ip        string
	closeOnce sync.Once
}
type chatGroupRoom struct {
	forward chan *model.ChatModel
	join    chan *chatGroupClient
	leave   chan *chatGroupClient
	clients map[*chatGroupClient]struct{}
}

var chatGroupUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var chatGroup = newChatGroupRoom()

func newChatGroupRoom() *chatGroupRoom {
	return &chatGroupRoom{
		forward: make(chan *model.ChatModel),
		join:    make(chan *chatGroupClient),
		leave:   make(chan *chatGroupClient),
		clients: make(map[*chatGroupClient]struct{}),
	}
}
func init() {
	go chatGroup.run()
}

func (c *chatGroupClient) trySend(msg []byte) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	select {
	case c.send <- msg:
		return true
	default:
		return false
	}
}

func (c *chatGroupClient) close() {
	c.closeOnce.Do(func() {
		c.room.leave <- c
		_ = c.conn.Close()
	})
}

func (c *chatGroupClient) readPump() {
	defer c.close()
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		p := model.ParamChatGroup{}
		if err := json.Unmarshal(raw, &p); err != nil {
			continue
		}
		p.NickName = c.nickName
		p.Avatar = c.avatar
		if p.MsgType == 0 {
			p.MsgType = ctype.TextMsg
		}
		cm := model.ChatModel{
			NickName: p.NickName,
			Avatar:   p.Avatar,
			Content:  p.Content,
			IP:       c.ip,
			Addr:     "",
			UserSN:   c.userSN,
			MsgType:  p.MsgType,
		}
		c.room.forward <- &cm
	}
}

func (c *chatGroupClient) writePump() {
	defer c.close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (r *chatGroupRoom) onlineCount() int {
	return len(r.clients)
}

func (r *chatGroupRoom) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = struct{}{}
			online := r.onlineCount()

			go func() {
				history, err := chat.ListRecentChatRecords(50)
				if err != nil {
					global.Log.Errorf("chat ListRecentChatRecords err:%s\n", err.Error())
					return
				}
				userSNSet := make(map[int64]struct{})
				for _, item := range history {
					if item.UserSN > 0 {
						userSNSet[item.UserSN] = struct{}{}
					}
				}
				userSNList := make([]int64, 0, len(userSNSet))
				for sn := range userSNSet {
					userSNList = append(userSNList, sn)
				}
				userAvatarMap := make(map[int64]string)
				if len(userSNList) > 0 {
					udl, err := user.GetUserDetailListBySNList(userSNList)
					if err == nil && udl != nil {
						for _, ud := range *udl {
							userAvatarMap[ud.SN] = ud.Avatar
						}
					}
				}
				for _, item := range history {
					avatar := item.Avatar
					if avatar == "" && item.UserSN > 0 {
						if a, ok := userAvatarMap[item.UserSN]; ok {
							avatar = a
						}
					}
					resp := model.ResponseChatGroup{
						ParamChatGroup: model.ParamChatGroup{
							NickName: item.NickName,
							Avatar:   avatar,
							Content:  item.Content,
							UserSN:   item.UserSN,
							MsgType:  item.MsgType,
						},
						Date:        item.CreateTime,
						OnlineCount: online,
					}
					b, err := json.Marshal(resp)
					if err != nil {
						continue
					}
					if ok := client.trySend(b); !ok {
						return
					}
				}
			}()

			r.broadcastGroupMsg(client.nickName+" 进入群聊", online, ctype.InRoomMsg)

		case client := <-r.leave:
			if _, ok := r.clients[client]; !ok {
				continue
			}
			delete(r.clients, client)
			close(client.send)
			r.broadcastGroupMsg(client.nickName+" 已退出群聊", r.onlineCount(), ctype.OutRoomMsg)

		case msg := <-r.forward:
			go func(cm *model.ChatModel) {
				if err := chat.CreateChatRecord(cm); err != nil {
					global.Log.Errorf("chat CreateChatRecord err:%s\n", err.Error())
				}
			}(msg)

			online := r.onlineCount()
			resp := model.ResponseChatGroup{
				ParamChatGroup: model.ParamChatGroup{
					NickName: msg.NickName,
					Avatar:   msg.Avatar,
					Content:  msg.Content,
					UserSN:   msg.UserSN,
					MsgType:  msg.MsgType,
				},
				Date:        time.Now(),
				OnlineCount: online,
			}
			b, err := json.Marshal(resp)
			if err != nil {
				continue
			}
			for client := range r.clients {
				if ok := client.trySend(b); !ok {
					delete(r.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (r *chatGroupRoom) broadcastGroupMsg(content string, online int, msgType ctype.MsgType) {
	resp := model.ResponseChatGroup{
		ParamChatGroup: model.ParamChatGroup{
			NickName: "System",
			Avatar:   "",
			Content:  content,
			MsgType:  msgType,
		},
		Date:        time.Now(),
		OnlineCount: online,
	}
	b, err := json.Marshal(resp)
	if err != nil {
		return
	}
	for client := range r.clients {
		if ok := client.trySend(b); !ok {
			delete(r.clients, client)
			close(client.send)
		}
	}
}

func HandleChatGroupWS(c *gin.Context) error {
	_claims, ok := c.Get("claims")
	if !ok || _claims == nil {
		return errors.New("need login")
	}

	conn, err := chatGroupUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return err
	}

	var nickName string
	var avatar string
	var userSN int64 = 0
	claims := _claims.(*jwt.MyClaims)
	userSN = claims.SN
	nickName = claims.Username
	ud, err := user.GetUserDetailBySN(claims.SN)
	if err == nil {
		avatar = ud.Avatar
	}
	if userSN <= 0 {
		_ = conn.Close()
		return errors.New("need login")
	}

	client := &chatGroupClient{
		conn:     conn,
		send:     make(chan []byte, 256),
		room:     chatGroup,
		nickName: nickName,
		avatar:   avatar,
		ip:       c.ClientIP(),
		userSN:   userSN,
	}
	chatGroup.join <- client
	go client.writePump()
	go client.readPump()
	return nil
}

func GetChatGroupRecords(p *model.ParamList) (*model.PageResult[model.ChatModel], error) {
	list, count, err := chat.GetChatRecords(p.Page, p.Size)
	if err != nil {
		return nil, err
	}
	return &model.PageResult[model.ChatModel]{
		List:  list,
		Count: count,
	}, nil
}

func UploadChatGroupImage(c *gin.Context, fh *multipart.FileHeader) (string, error) {
	if fh == nil {
		return "", errors.New("image required")
	}
	fileExt := strings.Split(fh.Filename, ".")
	if len(fileExt) != 2 {
		return "", errors.New("invalid image")
	}
	ext := strings.ToLower(fileExt[1])
	if _, exists := model.WhiteImageExtList[ext]; !exists {
		return "", errors.New("invalid image")
	}
	size := float64(fh.Size) / 1024 / 1024
	if size >= float64(global.Config.Upload.Size) {
		return "", errors.New("image too large")
	}
	chatDir := filepath.Join(global.Config.Upload.Path, "chat")
	file.CreateFolder(chatDir)
	targetPath := filepath.Join(chatDir, fh.Filename)
	if err := c.SaveUploadedFile(fh, targetPath); err != nil {
		return "", err
	}
	return "/uploads/file/chat/" + fh.Filename, nil
}
