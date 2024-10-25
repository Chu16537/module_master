package mwebscoketserver

import (
	"time"

	"github.com/gorilla/websocket"
)

type IClient interface {
	GetUid() int64                 // 取得id
	SetLastReadMsgTime(unix int64) // 設定最後訊息時間
	GetLastReadMsgTime() int64     // 取得最後訊息時間
	WriteMessageQueue(msg []byte)  // 傳送資料 順序
	WriteMessage(msg []byte)       // 傳送資料 不管順序
	Done()                         // 斷線
	WaitDone()                     // 等待斷線
}

type Client struct {
	conn            *websocket.Conn
	uid             int64
	doneChan        chan bool
	sender          chan []byte
	lastReadMsgTime int64
}

func newClient(conn *websocket.Conn, uid int64, sender chan []byte) IClient {
	return &Client{
		conn:            conn,
		uid:             uid,
		doneChan:        make(chan bool, 1),
		sender:          sender,
		lastReadMsgTime: time.Now().Unix(),
	}
}

func (c *Client) GetUid() int64 {
	return c.uid
}

func (c *Client) SetLastReadMsgTime(unix int64) {
	c.lastReadMsgTime = unix
}

func (c *Client) GetLastReadMsgTime() int64 {
	return c.lastReadMsgTime
}

func (c *Client) WriteMessageQueue(msg []byte) {
	c.sender <- msg
}

func (c *Client) WriteMessage(msg []byte) {
	c.conn.WriteMessage(websocket.TextMessage, msg)
}

func (c *Client) Done() {
	c.doneChan <- true
}

func (c *Client) WaitDone() {
	<-c.doneChan
}
