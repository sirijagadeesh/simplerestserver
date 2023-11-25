package main

import (
	"flag"

	"github.com/sirijagadeesh/simplerestserver/api"
)

func main() {
	port := 0
	flag.IntVar(&port, "listen-addr", 3008, "server listen port")
	flag.Parse()

	api.NewServer(port).Start()
}
