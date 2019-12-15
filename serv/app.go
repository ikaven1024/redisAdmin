package main

import (
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"github/ikaven/redisAdmin/web"

	"github/ikaven/redisAdmin/api"
	"github/ikaven/redisAdmin/config"
	"github/ikaven/redisAdmin/repository"
	"github/ikaven/redisAdmin/server"
)

type App struct {
	config *config.Config

	mux *http.ServeMux
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
	serverManager := server.NewManager(db)

	app.mux = http.NewServeMux()
	app.mux.Handle("/", web.Server())
	app.mux.Handle("/api", api.SetupApi(serverManager))

	return &app
}

func (a *App) Run() {
	addr := ":80"
	if len(a.config.Server.Address) != 0 {
		addr = a.config.Server.Address
	}

	log.Println("Server listening at:", addr)
	log.Println(http.ListenAndServe(addr, a.mux))
}
