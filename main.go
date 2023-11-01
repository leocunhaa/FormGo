package main

import (
	"log"
	"module/handlers"
	"net/http"
)

func main() {
	http.HandleFunc("/", handlers.SubscriptionHandler)

	if err := http.ListenAndServe(":4040", nil); err != nil {
		log.Fatal(err)
	}
}
