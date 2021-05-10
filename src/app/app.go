package app

import (
	"framework/cfgargs"
	"framework/db"
	"framework/logger"
	"framework/net/socketio"
	"gate/handler"
	sio "github.com/googollee/go-socket.io"
	"sync"
)

var (
	once sync.Once
	app  *App
)

type App struct {
	sioSrv *socketio.Server
	cfg    *cfgargs.SrvConfig
}

func GetApp() *App {
	once.Do(func() {
		a := &App{}
		app = a
	})
	return app
}

func (a *App) Init(cfg *cfgargs.SrvConfig) {
	// db
	db.InitRedisClient(cfg)
	err := db.InitMongoClient(cfg)
	if err != nil {
		logger.Fatal("init mongo db err: %v", err)
		return
	}
	a.InitSocketIO(cfg)

}

func (a *App) InitSocketIO(cfg *cfgargs.SrvConfig) {
	// sio srv init
	a.sioSrv = socketio.NewServer()

	a.sioSrv.OnConnect(func(conn sio.Conn) error {
		logger.Info("socket.io connected, socket id :%v", conn.ID())
		return nil
	})

	a.sioSrv.OnDisconnect(func(conn sio.Conn, s string) {
		logger.Info("socket.io disconnected, socket id :%v", conn.ID())
	})

	a.sioSrv.OnError(func(conn sio.Conn, err error) {
		logger.Error("socket.io on err: %v, id: %v", err, conn.ID())
	})

	// mount handlers
	handlers := socketio.NewSIOHandlers()
	handlers["hello"] = handler.Hello
	handlers["chat"] = handler.Chat
	a.sioSrv.MountHandlers("/", handlers)

	// run
	go a.sioSrv.Run(cfg) //nolint: errcheck
}
