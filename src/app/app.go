package app

import (
	"framework/cfgargs"
	"gate/server"
	"sync"
)

var (
	once sync.Once
	app  *App
)

type App struct {
	srv *server.Server
	cfg *cfgargs.SrvConfig
}

func GetApp() *App {
	once.Do(func() {
		a := &App{}
		app = a
	})
	return app
}

func (a *App) Init(cfg *cfgargs.SrvConfig) {

	//socket.io
	a.srv = server.NewServer()
	a.srv.Init(cfg)
}
