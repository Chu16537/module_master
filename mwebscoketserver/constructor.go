package mwebscoketserver

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Chu16537/module_master/errorcode"
	"github.com/Chu16537/module_master/mjson"
	"github.com/Chu16537/module_master/mlog"
	"github.com/Chu16537/module_master/mtime"
	"github.com/Chu16537/module_master/muid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type Config struct {
	Addr               string
	MaxConn            int
	AliveTimeoutSecond int
}

type Handler struct {
	ctx      context.Context
	config   *Config
	log      mlog.ILog
	ws       http.Server
	upgrader websocket.Upgrader
	ih       IHandler

	uid       *muid.Handler
	lock      sync.RWMutex
	clientMap map[int64]IClient // map[唯一編號]IClient
}

type IHandler interface {
	// 收到使用者訊號
	ReadMessage(*ClientReq)
}

var h *Handler

var upgrader = websocket.Upgrader{
	// WebSocket请求来自任何来源
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func New(ctx context.Context, config *Config, uid *muid.Handler, log mlog.ILog, ih IHandler) error {
	// 基本判斷
	if config.Addr == "" {
		return errors.New("websocket new error addr nil")
	}

	h = &Handler{
		ctx:       ctx,
		config:    config,
		log:       log,
		upgrader:  upgrader,
		ih:        ih,
		uid:       uid,
		lock:      sync.RWMutex{},
		clientMap: make(map[int64]IClient),
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
			opt := &mlog.LogData{
				Data: "mwebscoketserver done",
			}
			h.log.Info(opt)
			return
		case <-ticker.C:
			nowUnix := time.Now().Unix()
			h.lock.RLocker().Lock()
			for _, client := range h.clientMap {
				// 最新的請求時間更後面
				if nowUnix < client.GetLastReadMsgTime() {
					continue
				}

				// 一段時間內沒有送請求
				if nowUnix-client.GetLastReadMsgTime() > int64(h.config.AliveTimeoutSecond) {
					// 斷線
					client.Done()
				}
			}

			h.lock.RLocker().Unlock()
		}
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 将 HTTP 连接升级为 WebSocket 连接
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		opt := &mlog.LogData{
			Err: errorcode.Server(err),
		}
		h.log.Error(opt)
		return
	}

	sendChan := make(chan []byte, 32)
	clientId := h.uid.CreateID()

	defer func() {
		conn.Close()

		h.lock.Lock()
		delete(h.clientMap, clientId)
		h.lock.Unlock()
	}()

	client := newClient(conn, clientId, sendChan)

	h.lock.Lock()
	h.clientMap[clientId] = client
	h.lock.Unlock()

	// 讀取訊息
	go h.reading(conn, client)

	// 傳送訊息
	go h.sending(conn, sendChan)

	// 斷線
	client.WaitDone()
}

func (h *Handler) reading(conn *websocket.Conn, client IClient) {
	defer func() {
		client.Done()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			opt := &mlog.LogData{
				Err: errorcode.Server(err),
			}
			h.log.Error(opt)
			return
		}

		// 格式轉換 有錯誤當作他亂打請求
		req := &ClientReq{}
		err = mjson.Unmarshal(msg, req)
		if err != nil {
			client.WriteMessage([]byte(reqUnmarshalErr))
			opt := &mlog.LogData{
				Err: errorcode.Server(err),
			}
			h.log.Warn(opt)
			continue
		}

		req.NewID(client.GetUid())

		// 更新最後請求時間
		client.SetLastReadMsgTime(mtime.GetZero().Unix())

		// 請求實作
		h.ih.ReadMessage(req)
	}

	fmt.Println("asd")
}

func (h *Handler) sending(conn *websocket.Conn, sender <-chan []byte) {
	for msg := range sender {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			opt := &mlog.LogData{
				Err: errorcode.Server(err),
			}
			h.log.Error(opt)
		}
	}
}

// 返回請求資料
func (h *Handler) Response(res *ClientRes) {
	reqId, clientId := res.GetId()

	c, ok := h.clientMap[clientId]

	if !ok {
		// 代表使用者斷線
		opt := &mlog.LogData{
			Err: errorcode.Server(errors.Errorf("clientId:%v is disconnect", clientId)),
		}
		h.log.Warn(opt)
		return
	}

	res.Id = reqId
	resByte, err := mjson.Marshal(res)
	if err != nil {
		opt := &mlog.LogData{
			Err: errorcode.Server(errors.Errorf("clientId:%v res marshal err data:%v", clientId, res)),
		}
		h.log.Error(opt)
		return
	}

	c.WriteMessageQueue(resByte)
}
