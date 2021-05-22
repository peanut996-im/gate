package server

import (
	"framework/api"
	"framework/logger"
	sio "github.com/googollee/go-socket.io"
)

const (
	ChatEvent         = "chat"
	AddFriendEvent    = "addFriend"
	DeleteFriendEvent = "deleteFriend"
	CreateGroupEvent  = "createGroup"
	JoinGroupEvent    = "joinGroup"
	LeaveGroupEvent   = "leaveGroup"
)

func Chat(conn sio.Conn) {
	logger.Info("/chat: revived message from %v", conn.ID())
}

//func Ping(conn sio.Conn) {
//	logger.Info("/pingpong from [%v]", conn.ID())
//	conn.Emit("pingpong", "pong")
//}
//
//func Console(conn sio.Conn, msg string) {
//	logger.Info("/console [%v]: %v", conn.ID(), msg)
//}

func AddFriend(conn sio.Conn, request *api.FriendRequest) {
	logger.Info("/addFriend from[%v]: %+v", conn.ID(), request)
}

func DeleteFriend(conn sio.Conn, request *api.FriendRequest) {
	logger.Info("/addFriend from[%v]: %+v", conn.ID(), request)
}

func CreateGroup(conn sio.Conn) {
	logger.Info("/createGroup from[%v]: %+v", conn.ID())
}

func JoinGroup(conn sio.Conn) {
	logger.Info("/joinGroup from[%v]: %+v", conn.ID())
}

func LeaveGroup(conn sio.Conn) {
	logger.Info("/leaveGroup from[%v]: %+v", conn.ID())
}
