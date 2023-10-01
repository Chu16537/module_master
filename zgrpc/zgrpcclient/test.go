package zgrpcclient

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Chu16537/gomodule/zgrpc/test"
	"google.golang.org/grpc"
)

func Test() {
	cxt := context.Background()
	conf := &Config{
		Addr: "127.0.0.1:50051",
	}

	h, err := New(cxt, conf)

	if err != nil {
		fmt.Println("err", err)
		return
	}
	readChan := make(chan *test.Res, 1024)
	sendChan := make(chan *test.Req, 1024)
	s := NewClient(h.conn, readChan, sendChan)

	go s.Run(cxt)

	time.Sleep(2 * time.Second)
	sendChan <- &test.Req{
		ReqID: 111,
	}
}

type Client struct {
	client   test.RoomserverClient
	stream   test.Roomserver_PingClient
	errChan  chan struct{}
	readChan chan *test.Res
	sendChan chan *test.Req
}

func NewClient(conn *grpc.ClientConn, readChan chan *test.Res, sendChan chan *test.Req) *Client {
	return &Client{
		client:   test.NewRoomserverClient(conn),
		readChan: readChan,
		sendChan: sendChan,
	}
}

// OperateRoom
func (c *Client) Run(cxt context.Context) bool {
	ctx, cancel := context.WithTimeout(cxt, 10*time.Second)
	defer cancel()

	// 開啟 channel 接收訊號
	c.errChan = make(chan struct{})

	stream, err := c.client.Ping(ctx)

	if err != nil {
		fmt.Println("stream fail", err)
		return false
	}

	c.stream = stream

	// read
	go c.read()

	// send
	go c.send()

	<-c.errChan
	fmt.Println("Run 接收到錯誤")
	return false
}

// 讀取
func (c *Client) read() {
	for {
		res, err := c.stream.Recv()

		if err != nil {
			fmt.Println("read stream recv err", err)
			c.errChan <- struct{}{}
			return
		}

		// 接收到訊號後實作
		fmt.Println("read res", res)
		// 把資料丟入實作層
		c.readChan <- res
	}
}

// 傳送
func (c *Client) send() {
	for {
		data := <-c.sendChan
		err := c.stream.Send(data)

		if err != nil {
			fmt.Println("grpc send fail", err)
		}
		if err == io.EOF {
			return
		}
	}
}
