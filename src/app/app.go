package app

import (
	"framework/cfgargs"
	"framework/db"
	"framework/logger"
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
	// db
	db.InitRedisClient(cfg)
	err := db.InitMongoClient(cfg)
	if err != nil {
		logger.Fatal("init mongo db err: %v", err)
		return
	}
	a.srv = server.NewServer()
	a.srv.InitSocketIO(cfg)
}
