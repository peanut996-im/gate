package handler

import (
	"fmt"
	"framework/logger"
	sio "github.com/googollee/go-socket.io"
)

func Hello(conn sio.Conn) interface{} {
	fmt.Println(conn.URL())
	fmt.Println("get hello")
	return "OK"
}

func Chat(conn sio.Conn, message string) interface{} {
	logger.Info("get message",fmt.Sprintf("%v: %v",conn.ID(),message))
	conn.Emit("chat",fmt.Sprintf("%v: %v",conn.ID(),message));
	return "ok"
}
