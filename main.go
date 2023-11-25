package main

import (
	"flag"

	"github.com/sirijagadeesh/simplerestserver/rest"
)

func main() {
	port := 0
	flag.IntVar(&port, "port", 3008, "server listen port")
	flag.Parse()

	rest.NewServer(port).Start()
}
