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

func (s *Server) Chat(conn sio.Conn) {
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

func (s *Server) AddFriend(conn sio.Conn, request *api.FriendRequest) {
	logger.Info("/addFriend from[%v]: %+v", conn.ID(), request)
}

func (s *Server) DeleteFriend(conn sio.Conn, request *api.FriendRequest) {
	logger.Info("/deleteFriend from[%v]: %+v", conn.ID(), request)
}

func (s *Server) CreateGroup(conn sio.Conn) {
	logger.Info("/createGroup from[%v]: ", conn.ID())
}

func (s *Server) JoinGroup(conn sio.Conn) {
	logger.Info("/joinGroup from[%v]: ", conn.ID())
}

func (s *Server) LeaveGroup(conn sio.Conn) {
	logger.Info("/leaveGroup from[%v]: ", conn.ID())
}
