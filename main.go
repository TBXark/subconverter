package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/TBXark/subconverter/api"
)

var BuildVersion = "dev"

func main() {
	addr := flag.String("addr", ":3000", "http listen address")
	help := flag.Bool("help", false, "show help")
	version := flag.Bool("version", false, "show version")

	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	if *version {
		fmt.Println(BuildVersion)
		return
	}

	engine := api.NewEngine()
	log.Printf("Starting server on %s", *addr)
	err := engine.Run(*addr)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Failed to start server: %v", err)
	}
}
