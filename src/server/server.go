package server

import (
	"encoding/json"
	"fmt"
	"framework/api"
	"framework/api/model"
	"framework/broker/logic"
	"framework/cfgargs"
	"framework/logger"
	"framework/tool"
	sio "github.com/googollee/go-socket.io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Server struct {
	srv                *sio.Server
	logicBroker        logic.LogicBroker
	nsp                string
	handlers           map[string]interface{}
	SocketIOToSessions map[string]*Session
	//UIDSceneToSessions map[string]*Session
	SceneToSessions map[string]*Session
	sync.Mutex
}

func NewSIOHandlers() map[string]interface{} {
	return make(map[string]interface{})
}

func NewServer() *Server {
	s := &Server{
		srv:                sio.NewServer(nil),
		nsp:                "/",
		handlers:           make(map[string]interface{}),
		SocketIOToSessions: make(map[string]*Session),
		//UIDSceneToSessions: make(map[string]*Session),
		SceneToSessions: make(map[string]*Session),
	}
	return s
}

func (s *Server) GetEventHandler(event string) interface{} {
	return func(conn sio.Conn, data interface{}) {
		logger.Info("/%v from[%v]: %+v", event, conn.ID(), data)
		rawJson, err := s.logicBroker.Send(event, data)
		if nil != err {
			conn.Emit(event, api.NewHttpInnerErrorResponse(err))
			logger.Error("Gate.Event[%v] Broker err: %v", event, err)
		}
		conn.Emit(event, rawJson.(json.RawMessage))
	}
}

func (s *Server) MountHandlers() {
	s.handlers[api.EventAddFriend] = s.GetEventHandler(api.EventAddFriend)
	s.handlers[api.EventJoinGroup] = s.GetEventHandler(api.EventJoinGroup)
	s.handlers[api.EventDeleteFriend] = s.GetEventHandler(api.EventJoinGroup)
	s.handlers[api.EventCreateGroup] = s.GetEventHandler(api.EventCreateGroup)
	s.handlers[api.EventChat] = s.GetEventHandler(api.EventChat)
	s.handlers[api.EventLeaveGroup] = s.GetEventHandler(api.EventLeaveGroup)

	for k, v := range s.handlers {
		s.srv.OnEvent(s.nsp, k, v)
	}
}

func (s *Server) Init(cfg *cfgargs.SrvConfig) {
	// rpc by http
	if cfg.Logic.Mode == "http" {
		s.logicBroker = logic.NewLogicBrokerHttp()
	}
	s.logicBroker.Init(cfg)
	// sio srv init

	s.OnConnect(func(conn sio.Conn) error {
		logger.Info("socket.io connected, socket id :%v", conn.ID())
		si := NewSession(conn)
		err := s.AcceptSession(si)

		if err != nil{
			go func() {
				//Reconnect time
				conn.Emit("auth", api.AuthFaildResp)
				<-time.After(20 * time.Second)
				conn.Close()
			}()
		} else {
			conn.Emit("auth", api.NewSuccessResponse(nil))
			go s.PushInitData(si)
		}
		return nil
	})

	s.OnDisconnect(func(conn sio.Conn, message string) {
		logger.Info("socket.io disconnected, socket id :%v", conn.ID())
		s.Lock()
		si, _ := s.SocketIOToSessions[conn.ID()]
		if nil != si {
			logger.Info("Session disconnected, session[%v]", si.ToString())
			delete(s.SocketIOToSessions, conn.ID())
		}

		s.Unlock()
		s.DisconnectSession(conn)
	})

	s.OnError(func(conn sio.Conn, err error) {
		logger.Error("socket.io on err: %v, id: %v", err, conn.ID())
	})

	s.MountHandlers()
	// run
	go s.Run(cfg) //nolint: errcheck
}

func (s *Server) Run(cfg *cfgargs.SrvConfig) error {
	defer func(srv *sio.Server) {
		err := srv.Close()
		if err != nil {
			panic(err)
		}
	}(s.srv)
	go func() {
		err := s.srv.Serve()
		if err != nil {
			panic(err)
		}
	}() //nolint: errcheck

	if cfg.SocketIO.Cors {
		http.HandleFunc("/socket.io/", func(w http.ResponseWriter, r *http.Request) {
			allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Headers", allowHeaders)
			}
			if r.Method == "OPTIONS" {
				return
			}
			r.Header.Del("Origin")
			s.srv.ServeHTTP(w, r)
		})
	} else {
		http.Handle("/socket.io/", s.srv)
	}

	addr := fmt.Sprintf(":%v", cfg.SocketIO.Port)
	logger.Info("Serving at %v...", addr)

	err := http.ListenAndServe(addr, nil)
	logger.Fatal("Serving at %v... err:%v", addr, err)
	return err
}

func (s *Server) OnConnect(f func(sio.Conn) error) {
	s.srv.OnConnect(s.nsp, f)
}

func (s *Server) OnDisconnect(f func(sio.Conn, string)) {
	s.srv.OnDisconnect(s.nsp, f)
}

func (s *Server) OnError(f func(sio.Conn, error)) {
	s.srv.OnError(s.nsp, f)
}

func (s *Server) SocketIOToSession(c sio.Conn) *Session {
	s.Lock()
	si, ok := s.SocketIOToSessions[c.ID()]
	s.Unlock()
	if !ok {
		logger.Warn("session not found")
		return nil
	}
	return si
}

//
//func (s *Server) UIDSceneToSession(uidScene string) *Session {
//	s.Lock()
//	si, ok := s.UIDSceneToSessions[uidScene]
//	s.Unlock()
//	if !ok {
//		logger.Warn("session not found")
//		return nil
//	}
//	return si
//}

//SetNameSpace reset namespace
func (s *Server) SetNameSpace(nsp string) {
	s.nsp = nsp
}

//AcceptSession authentication for session
func (s *Server) AcceptSession(session *Session) error {
	ok,err := s.Auth(session)
	if !ok || nil != err{
		return err
	}
	s.Lock()
	s.SocketIOToSessions[session.GetID()] = session
	s.SceneToSessions[session.scene] = session
	s.Unlock()
	logger.Info("Session.Accept succeed, session:[%v]", session.ToString())

	logger.Info("Session.Accept done. Session[%v]", session.ToString())
	return nil
}

func (s *Server) DisconnectSession(conn sio.Conn) *Session {

	s.Lock()
	si, ok := s.SocketIOToSessions[conn.ID()]
	if ok || nil != si {
		delete(s.SocketIOToSessions, si.Conn.ID())
	} else {
		logger.Warn("Sessions.DisconnectSession[%v] not found", ToString(conn))
	}

	if nil != si {
		siScene, ok := s.SceneToSessions[si.scene]
		if ok || nil != siScene {
			logger.Info("Sessions.DisconnectSession,UIDAndScene:v%", si.scene)
			delete(s.SceneToSessions, si.scene)
		}
	}

	s.Unlock()
	return si
}


func (s *Server) Auth(session *Session) (bool, error) {
	vals, err := url.ParseQuery(session.query)
	sign, _ := api.MakeSignWithQueryParams(vals, cfgargs.GetLastSrvConfig().AppKey)
	if sign != vals.Get("sign") {
		logger.Info("Session.Auth failed. sign invalid: %v", sign)
		return false,api.ErrorCodeToError(api.ErrorSignInvalid)
	}
	if nil != err {
		logger.Info("parse token failed, err: %v", err)
		return false,api.ErrorCodeToError(api.ErrorTokenInvalid)
	}

	t := vals.Get("token")
	rawJson, err := s.logicBroker.Send(api.EventAuth, t)
	if err != nil {
		logger.Error("Session.Auth get auth response err. err: %v", err)
		return false,api.ErrorCodeToError(api.ErrorHttpInnerError)
	}
	resp := &api.BaseRepsonse{}
	if err = json.Unmarshal(rawJson.(json.RawMessage), resp); err != nil {
		logger.Info("Session.Save json unmarshal err. err:%v, Session:[%v]", err, session.ToString())
		return false, api.ErrorCodeToError(api.ErrorHttpInnerError)
	}
	if resp.Code != api.ErrorCodeOK || resp.Data == nil{
		// Auth failed
		logger.Error("Session.Auth auth failed. Maybe token expired or user not exist? UID: [%v], Session:[%v]", vals.Get("uid"), session.ToString())
		return false,api.ErrorCodeToError(api.ErrorAuthFailed)
	}
	u := &model.User{}
	if err = tool.MapToStruct(resp.Data, u); err != nil {
		logger.Info("Session.Auth json unmarshal err. err:%v, Session:[%v]", err, session.ToString())
		return false,api.ErrorCodeToError(api.ErrorHttpInnerError)
	}

	logger.Info("Session.Auth succeed.")
	session.SetScene(u.UID)
	session.token = t
	return true,nil
}

// PushInitData PushLoadData push init data
func (s *Server) PushInitData(si *Session) {
	data, err := s.logicBroker.Send(api.EventLoad, si.GetScene())
	if err != nil {
		return
	}
	si.Push("load", data)
}
