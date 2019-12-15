package web

import (
	"flag"
	"net/http"
)

var webRoot = flag.String("web-root", "./wwww", "dir of web resource.")

func Server() http.Handler {
	return http.FileServer(http.Dir(*webRoot))
}
