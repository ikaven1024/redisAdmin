package server

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/ikaven1024/redisAdmin/api"
	"github.com/ikaven1024/redisAdmin/redis_server"
)

var webRoot = flag.String("web-root", "./wwww", "dir of web resource.")

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	gin.ForceConsoleColor()

	engine := gin.Default()
	engine.StaticFile(".", *webRoot)
	engine.StaticFile("/index.html", filepath.Join(*webRoot, "index.html"))
	engine.StaticFile("/favicon.ico", filepath.Join(*webRoot, "favicon.ico"))
	engine.Static("/static", filepath.Join(*webRoot, "static"))

	return &Server{
		engine: engine,
	}
}

func (s *Server) InstallApi(serverManager *redis_server.Manager) {
	api.InstallApi(s.engine, serverManager)
}

func (s *Server) Run(addr string) {
	log.Println("Server listening at:", addr)
	log.Println(http.ListenAndServe(addr, s.engine))
}
