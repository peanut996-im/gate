package server

import (
	"fmt"
	"framework/api"
	"framework/api/model"
	"framework/logger"

	sio "github.com/googollee/go-socket.io"
)

type Session struct {
	Conn  sio.Conn
	token string
	scene string
}

func NewSession(conn sio.Conn) *Session {
	return &Session{
		Conn: conn,
	}
}
func (s *Session) GetScene() string {
	return s.scene
}

func (s *Session) SetScene(scene string) {
	s.scene = scene
}

func (s *Session) GetID() string {
	return s.Conn.ID()
}

func ToString(c sio.Conn) string {
	if nil != c {
		id := c.ID()
		localAddr := c.LocalAddr()
		remoteAddr := c.RemoteAddr()
		return fmt.Sprintf("ID:%v addr.local:%v addr.remote:%v", id, localAddr, remoteAddr)
	}
	return "conn not found"
}

//func (s *Session) UIDSceneString() string {
//	return fmt.Sprintf("uid:%v_scene:%v", s.uid, s.scene)
//}

func (s *Session) Auth(token string) (*model.User, error) {
	//todo move auth to logic
	resp, err := api.CheckToken(token)
	if err != nil {
		logger.Info("check user token failed error: %v", err)
		return nil, err
	}
	return resp, nil
}

func (s *Session) ToString() string {
	return ToString(s.Conn)
}

func (s *Session) Push(event string, data interface{}) {
	s.Conn.Emit(event, data)
}
