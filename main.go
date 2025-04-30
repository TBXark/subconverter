package main

import (
	"log"
	"net/http"

	"github.com/TBXark/subconverter/api"
)

func main() {
	http.HandleFunc("/", api.Handler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
