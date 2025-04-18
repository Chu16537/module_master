package mwebscoketserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mjson"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Config struct {
	Port               string
	MaxConn            int
	AliveTimeoutSecond int
}

type Handler struct {
	ctx            context.Context
	config         *Config
	ws             http.Server
	upgrader       websocket.Upgrader
	lock           sync.RWMutex
	clientConnents []IClient // map[唯一編號]IClient
	ih             IHandler
}

type IHandler interface {
	// 連線
	Connect(client IClient) error
	// 斷線
	Disconnect(idx uint32)
	// 傳遞資料
	ReadMessage(*ToHanglerReq)
}

var h *Handler

var upgrader = websocket.Upgrader{
	// WebSocket请求来自任何来源
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(ctx context.Context, config *Config, ih IHandler) error {
	// 基本判斷
	if config.Port == "" {
		return errors.New("websocket new error addr nil")
	}

	h = &Handler{
		ctx:            ctx,
		config:         config,
		upgrader:       upgrader,
		lock:           sync.RWMutex{},
		clientConnents: make([]IClient, config.MaxConn),
		ih:             ih,
	}

	h.ws = http.Server{Addr: fmt.Sprintf(":%v", h.config.Port), Handler: h}

	go checkAlive()

	errChan := make(chan error, 1)
	go func() {
		err := h.ws.ListenAndServe()

		if err != nil {
			errChan <- err
		}
	}()

	// 等待n秒判斷是否有錯
	select {
	case err := <-errChan:
		return err

	case <-time.After(5 * time.Second):
		// 等待5秒發現沒有錯誤
		fmt.Println("websocket new success max conn:", h.config.MaxConn)
		return nil
	}

}

// 存活檢查
func checkAlive() {
	// 创建一个新的 ticker
	ticker := time.NewTicker(time.Duration(h.config.AliveTimeoutSecond) * time.Second)
	defer ticker.Stop() // 在退出时停止 ticker

	for {
		select {
		case <-h.ctx.Done():
			return

		case <-ticker.C:
			nowUnix := time.Now().Unix()

			h.lock.Lock()
			for _, client := range h.clientConnents {
				if client == nil {
					continue
				}

				// 一段時間內沒有送請求
				if nowUnix-client.GetLastReadMsgTime() > int64(h.config.AliveTimeoutSecond) {
					// 斷開客戶端連線
					h.deleteUser(client)
				}
			}

			h.lock.Unlock()
		}
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer conn.Close()

	// 創建client
	client, err := h.createUser(conn)
	if err != nil {
		fmt.Println("createUser", err)
		return
	}

	// 讀取訊息
	go h.reading(client)

	// 傳送訊息
	go h.sending(client)

	// 等待斷線
	client.WaitDone()
	// 刪除玩家
	h.deleteUser(client)
}

func (h *Handler) reading(client IClient) {
	conn := client.GetConn()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// 格式轉換 有錯誤
		req := &ClientReq{}
		err = mjson.Unmarshal(msg, req)
		if err != nil {
			er := errorcode.New(errorcode.Code_Data_Unmarshal_Error, err)
			b, _ := json.Marshal(er)
			conn.WriteMessage(websocket.TextMessage, b)
			continue
		}

		toHanglerReq := &ToHanglerReq{
			ClientId: client.GetUid(),
			Req:      req,
		}

		// 更新最後請求時間
		client.UpdateLastReadTime(time.Now().Unix())

		// 請求實作
		h.ih.ReadMessage(toHanglerReq)
	}
}

func (h *Handler) sending(client IClient) {
	conn := client.GetConn()
	sender := client.GetSender()

	for msg := range sender {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

// func (h *Handler) clientDone(client IClient) {
// 	// 通知實做層
// 	h.ih.Disconnect(client.GetUid())
// 	// 斷線
// 	client.Done()
// }

// 返回請求資料
func Response(res *ToHanglerRes) error {
	if res == nil {
		return fmt.Errorf("Response res nil")
	}

	c := h.clientConnents[res.ClientId]

	// 代表連線已經斷線
	if c == nil {
		return nil
	}

	// 回傳前端
	resByte, err := mjson.Marshal(res.Res)
	if err != nil {
		return err
	}

	c.WriteMessageQueue(resByte)
	return nil
}

// 創建user
func (h *Handler) createUser(conn *websocket.Conn) (IClient, error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	id := -1

	for i, v := range h.clientConnents {
		if v == nil {
			id = i
			break
		}
	}

	if id == -1 {
		return nil, fmt.Errorf("clientConnents is full")
	}

	sendChan := make(chan []byte, 128)
	client := newClient(conn, uint32(id), sendChan)

	// 創建錯誤
	err := h.ih.Connect(client)
	if err != nil {
		return nil, err
	}

	h.clientConnents[id] = client

	return client, nil
}

func (h *Handler) deleteUser(client IClient) {
	// 通知實做層
	h.ih.Disconnect(client.GetUid())
	// 斷線
	client.Done()

	h.lock.Lock()
	defer h.lock.Unlock()

	if h.clientConnents[client.GetUid()] != nil {
		h.clientConnents[client.GetUid()] = nil
	}
}
