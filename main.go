package main

import (
	"fmt"
	"net/http"
	"os"
	"smart-software-engineering/rate-calculator/web"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func main() {
	/*
		store, err := memory.NewStore()
		if err != nil {
			log.Fatal(err)
		}

		sessions, err := web.NewSessionManager()
		if err != nil {
			log.Fatal(err)
		}

		csrfKey := []byte("01234567890123456789012345678901")
		h := web.NewHandler(store, sessions, csrfKey)
	*/

	h := web.NewHandler()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	address := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting server on %s\n", address)
	http.ListenAndServe(address, h)
}
