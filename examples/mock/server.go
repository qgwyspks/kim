package mock

import (
	"errors"
	"kim"
	"time"
)

type ServerDemo struct {
}

// demo入口方法
func (s *ServerDemo) Start(id, protocol, addr string) {
	var srv kim.Server
	service := &naming.DefaultService{
		Id:       id,
		Protocol: protocol,
	}

	// 忽略 NewServer 的内部逻辑，你可以认为它是一个空的方法，或者一个 mock 对象
	if protocol == "ws" {
		srv = websocket.NewServer(addr, service)
	} else if protocol == "tcp" {
		srv = tcp.NewServer(addr, service)
	}

	handler := &ServerHandler{}

	srv.SetReadWait(time.Minute)
	srv.SetAcceptor(handler)
	srv.SetMessageListener(handler)
	srv.SetStateListener(handler)
}

// ServerHandler ServerHandler
type ServerHandler struct {
}

// Accept this connection
func (h *ServerHandler) Accept(conn kim.Conn, timeout time.Duration) (string, error) {
	// 1. 读取：客户端发送的鉴权数据包
	frame, err := conn.ReadFrame()
	if err != nil {
		return "", err
	}
	// 2. 解析：数据包内容就是 userId
	userID := string(frame.GetPayload())
	// 3. 鉴权：这里只是为了示例做一个 fake 验证，非空
	if userID == "" {
		return "", errors.New("user id is invalid")
	}
	return userID, nil
}

// Receive default listener
func (h *ServerHandler) Receive(ag kim.Agent, payload []byte) {
	ack := string(payload) + "from server "
	_ = ag.Push([]byte(ack))
}

// Disconnect default listener
func (h *ServerHandler) Disconnect(id string) error {
	logger.Warnf("disconnect %s", id)
	return nil
}
