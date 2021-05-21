package handler

import (
	"fmt"
	"framework/api/model"
	"framework/logger"
	sio "github.com/googollee/go-socket.io"
)

func Chat(conn sio.Conn, message *model.ChatMessage) {
	logger.Info("/chat: revived message from [%v]: %v", conn.ID(), message)
	conn.Emit("chat", fmt.Sprintf("%v: %v", conn.ID(), message))
}

func Ping(conn sio.Conn) {
	logger.Info("/pingpong from [%v]", conn.ID())
	conn.Emit("pingpong", "pong")
}

func Console(conn sio.Conn, msg string) {
	logger.Info("/console [%v]: %v", conn.ID(), msg)
}
