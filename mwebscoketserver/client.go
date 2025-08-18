package mwebscoketserver

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type IClient interface {
	GetConn() *websocket.Conn
	GetUid() uint                  // 取得id
	UpdateLastReadTime(unix int64) // 設定最後訊息時間
	GetLastReadMsgTime() int64     // 取得最後訊息時間
	GetSender() chan []byte
	WriteMessageQueue(msg []byte)  // 傳送資料 順序
	WriteMessage(msg []byte) error // 傳送資料 不管順序
	Done()                         // 斷線
	WaitDone()                     // 等待斷線
}

type Client struct {
	closeOnce       sync.Once
	conn            *websocket.Conn
	uid             uint
	doneChan        chan bool
	sender          chan []byte
	lastReadMsgTime int64
}

func newClient(conn *websocket.Conn, uid uint, sender chan []byte) IClient {
	return &Client{
		closeOnce:       sync.Once{},
		conn:            conn,
		uid:             uid,
		doneChan:        make(chan bool, 1),
		sender:          sender,
		lastReadMsgTime: time.Now().Unix(),
	}
}

func (c *Client) GetConn() *websocket.Conn {
	return c.conn
}

func (c *Client) GetUid() uint {
	return c.uid
}

func (c *Client) UpdateLastReadTime(unix int64) {
	c.lastReadMsgTime = unix
}

func (c *Client) GetLastReadMsgTime() int64 {
	return c.lastReadMsgTime
}

func (c *Client) GetSender() chan []byte {
	return c.sender
}

func (c *Client) WriteMessageQueue(msg []byte) {
	c.sender <- msg
}

func (c *Client) WriteMessage(msg []byte) error {
	return c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Client) Done() {
	c.closeOnce.Do(func() {
		close(c.doneChan)
	})
}

func (c *Client) WaitDone() {
	<-c.doneChan
}
