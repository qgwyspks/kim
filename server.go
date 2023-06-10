package kim

import (
	"context"
	"net"
	"time"
)

type Server interface {
	SetAcceptor(Acceptor)               // 用于设置一个 Acceptor,在 Server 的 Start() 方法中监听到连接之后，就要调用这个 Accept 方法让上层业务处理握手相关工作
	SetMessageListener(MessageListener) // 用于设置一个消息监听器
	SetStateListener(StateListener)     // 设置一个状态监听器，将连接断开的事件上报给业务层，让业务层可以实现一些逻辑处理
	SetReadWait(duration time.Duration) // 用于设置连接读超时，用于控制心跳逻辑
	SetChannelMap(ChannelMap)           // 设置一个连接管理器，Server 在内部会自动管理连接的生命周期

	Start() error
	Push(string, []byte) error
	Shutdown(ctx context.Context) error
}

type Acceptor interface {
	Accept(Conn, time.Duration) (string, error)
}

type StateListener interface {
	Disconnect(string) error
}

type MessageListener interface {
	Receive(Agent, []byte)
}

type Agent interface {
	ID() string
	Push([]byte) error
}

// Conn Connection
type Conn interface {
	net.Conn
	ReadFrame() (Frame, error)
	WriteFrame(OpCode, []byte) error
	Flush() error
}

type Channel interface {
	Conn
	Agent
	Close() error
	Readloop(lst MessageListener) error
	SetWriteWait(time.Duration)
	SetReadWait(time.Duration)
}

// Client is interface of client side
type Client interface {
	ID() string
	Name() string
	Connect(string) error // 设置向一个服务器地址发起连接
	SetDialer(Dialer)     // 设置一个拨号器，这个方法会在 Connect 中被调用，完成连接的建立和握手
	Send([]byte) error    // 发送消息到服务端
	Read() (Frame, error) // 读取一帧数据，这里底层复用了 kim.Conn,所以直接返回 Frame
	Close()               // 断开连接，退出
}

type Dialer interface {
	DialAndHandshake(DialerContext) (net.Conn, error)
}

type DialerContext struct {
	Id      string
	Name    string
	Address string
	Timeout time.Duration
}

type OpCode byte

const (
	OpContinuation OpCode = 0x0
	OpText         OpCode = 0x1
	OpBinary       OpCode = 0x2
	OpClose        OpCode = 0x8
	OpPing         OpCode = 0x9
	OpPong         OpCode = 0xa
)

type Frame interface {
	SetOpCode(OpCode)
	GetOpCode() OpCode
	SetPayload([]byte)
	GetPayload() []byte
}
