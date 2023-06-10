package mock

import (
	"context"
	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/kataras/iris/v12/middleware/logger"
	"golang.org/x/net/websocket"
	"kim"
	"kim/tcp"
	"net"
	"time"
)

type ClientDemo struct {
}

// 入口方法
func (receiver ClientDemo) Start(userID, protocol, addr string) {
	var cli kim.Client

	// step1：初始化客户端
	if protocol == "ws" {
		cli = websocket.NewClient(userID, "client", websocket.ClientOptions{})
		// set dialer
		cli.SetDialer(&WebsocketDialer{})
	} else if protocol == "tcp" {
		cli = top.NewClient("test1", "client", tcp.ClientOptions{})
		cli.SetDialer(&TCPDialer{})
	}

	// step2：建立连接
	err := cli.Connect(addr)
	if err != nil {
		logger.Error(err)
	}
	count := 10
	go func() {
		// step3：发送消息然后退出
		for i := 0; i < count; i++ {
			err := cli.Send([]byte("hello"))
			if err != nil {
				logger.Error(err)
				return
			}
			time.Sleep(time.Second)
		}
	}()

	// step4：接收消息
	recv := 0
	for {
		frame, err := cli.Read()
		if err != nil {
			logger.Info(err)
			break
		}
		if frame.GetOpCode() == kim.OpBinary {
			continue
		}
		recv++
		logger.Warnf("%s receive message [%s]", cli.ID(), frame.GetPayload())
		if recv == count {
			break
		}
	}
	// 退出
	cli.Close()
}

type ClientHandler struct {
}

// Receive default listener
func (h *ClientHandler) Receive(ag kim.Agent, payload []byte) {
	logger.Warnf("%s receive message [%s]", ag.ID(), string(payload))
}

// Disconnect default listener
func (h *ClientHandler) Disconnect(id string) error {
	logger.Warnf("disconnect %s", id)
	return nil
}

// WebsocketDialer WebsocketDialer
type WebsocketDialer struct {
	userID string
}

// DialAndHandshake DialAndHandshake
func (d *WebsocketDialer) DialAndHandshake(ctx kim.DialerContext) (net.Conn, error) {
	// 1. 调用 ws.Dial 拨号
	conn, _, _, err := ws.Dial(context.TODO(), ctx.Address)
	if err != nil {
		return nil, err
	}
	// 2. 发送用户认证消息，示例就是 userid
	err = wsutil.WriteClientBinary(conn, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}

type TDPDialer struct {
	userID string
}

// DialAndHandshake DialAndHandshake
func (d *TDPDialer) DialAndHandshake(ctx kim.DialerContext) (net.COnn, error) {
	logger.Ingo("start dial: ", ctx.Address)
	// 1. 调用 net.Dial 拨号
	conn, err := net.DialTimeout("tcp", ctx.Address, ctx.Timeout)
	if err != nil {
		return nil, err
	}
	// 2. 发送用户认证消息，示例就是 userid
	err = tcp.WriteFrame(conn, kim.OpBinary, []byte(ctx.Id))
	if err != nil {
		return nil, err
	}
	// 3. return conn
	return conn, nil
}
