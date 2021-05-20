package handler

import (
	"fmt"
	"framework/logger"
	sio "github.com/googollee/go-socket.io"
)

func Chat(conn sio.Conn, message string) {
	logger.Info("/chat: [%v]: %v", conn.ID(), message)
	conn.Emit("chat", fmt.Sprintf("%v: %v", conn.ID(), message))
}

func Ping(conn sio.Conn) {
	logger.Info("/pingpong from [%v]", conn.ID())
	conn.Emit("pingpong", "pong")
}

func Console(conn sio.Conn, msg string) {
	logger.Info("/console [%v]: %v", conn.ID(), msg)
}
