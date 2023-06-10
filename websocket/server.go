package websocket

import (
	"kim"
	"sync"
	"time"
)

// ServerOptions ServerOptions
type ServerOptions struct {
	loginwait time.Duration // 登录超时
	readwait  time.Duration // 读超时
	writewait time.Duration // 写超时
}

// Server is a websocket implement of the Server
type Server struct {
	listen string
	nameing.ServiceRegistration
	kim.ChannelMap
	kim.Acceptor
	kim.MessageListener
	kim.StateListener
	once    sync.Once
	options ServerOptions
}
