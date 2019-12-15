package main

import (
	"flag"
)

func main() {
	flag.Parse()
	app := NewApp()
	app.Run()
}
