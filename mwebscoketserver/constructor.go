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
	Addr               string
	MaxConn            int
	AliveTimeoutSecond int
}

type Handler struct {
	ctx            context.Context
	config         *Config
	ws             http.Server
	upgrader       websocket.Upgrader
	lock           sync.RWMutex
	clientIdx      uint32
	clientConnents []IClient // map[唯一編號]IClient
	ih             IHandler
}

type IHandler interface {
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
	if config.Addr == "" {
		return errors.New("websocket new error addr nil")
	}

	h = &Handler{
		ctx:            ctx,
		config:         config,
		upgrader:       upgrader,
		lock:           sync.RWMutex{},
		clientIdx:      0,
		clientConnents: make([]IClient, config.MaxConn),
		ih:             ih,
	}

	h.ws = http.Server{Addr: h.config.Addr, Handler: h}

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
					client.Done() // 斷開客戶端連線
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

	var (
		id     uint32 = 0
		isJoin bool   = false
	)
	sendChan := make(chan []byte, 128)

	defer func() {
		close(sendChan)
		conn.Close()

		// 斷線刪除
		h.lock.Lock()
		h.clientConnents[id] = nil
		h.lock.Unlock()
	}()

	h.lock.Lock()
	for i, v := range h.clientConnents {
		if v == nil {
			isJoin = true
			id = uint32(i)
			break
		}
	}
	h.lock.Unlock()

	if !isJoin {
		fmt.Println("clientConnents is full")
		return
	}

	client := newClient(conn, id, sendChan)

	h.clientConnents[id] = client

	// 讀取訊息
	go h.reading(conn, client)

	// 傳送訊息
	go h.sending(conn, sendChan)

	// 斷線
	client.WaitDone()
}

func (h *Handler) reading(conn *websocket.Conn, client IClient) {
	defer func() {
		// 通知實做層
		h.ih.Disconnect(client.GetUid())
		// 斷線
		client.Done()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// 格式轉換 有錯誤
		req := &ToHanglerReq{}
		err = mjson.Unmarshal(msg, req)
		if err != nil {
			er := errorcode.New(errorcode.Code_Data_Unmarshal_Error, err)
			b, _ := json.Marshal(er)
			conn.WriteMessage(websocket.TextMessage, b)
			continue
		}

		req.ClientId = client.GetUid()

		// 更新最後請求時間
		client.UpdateLastReadTime(time.Now().Unix())

		// 請求實作
		h.ih.ReadMessage(req)
	}
}

func (h *Handler) sending(conn *websocket.Conn, sender <-chan []byte) {
	for msg := range sender {
		conn.WriteMessage(websocket.TextMessage, msg)
	}
}

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
	clientRes := &ClientRes{
		RequestId: res.RequestId,
		Data:      res.Data,
	}

	resByte, err := mjson.Marshal(clientRes)
	if err != nil {
		return err
	}

	c.WriteMessageQueue(resByte)
	return nil
}
