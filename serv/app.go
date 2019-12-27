package main

import (
	"log"

	"github.com/jinzhu/gorm"

	"github.com/ikaven1024/redisAdmin/config"
	"github.com/ikaven1024/redisAdmin/redis_server"
	"github.com/ikaven1024/redisAdmin/repository"
	"github.com/ikaven1024/redisAdmin/server"
)

type App struct {
	config *config.Config

	server *server.Server
}

func NewApp() *App {
	var app App
	var err error
	var db *gorm.DB

	if app.config, err = config.LoadConfig(); err != nil {
		log.Fatalf("Load config error: %v", err)
	}

	if db, err = repository.Open(); err != nil {
		log.Fatalf("Open repository error: %v", err)
	}
	serverManager := redis_server.NewManager(db)

	app.server = server.New()
	app.server.InstallApi(serverManager)

	return &app
}

func (a *App) Run() {
	addr := ":80"
	if len(a.config.Server.Address) != 0 {
		addr = a.config.Server.Address
	}
	a.server.Run(addr)
}
