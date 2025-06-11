package main

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"smart-software-engineering/rate-calculator/rates"
	"smart-software-engineering/rate-calculator/web"
)

//go:embed static/**
var StaticFS embed.FS

//go:embed templates
var TemplateFS embed.FS

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
	file, err := os.Open("data/schedules/romania.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	h := web.NewHandler(rates.NewRateCalculator(file), StaticFS, TemplateFS)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	address := fmt.Sprintf(":%s", port)
	fmt.Printf("Starting server on %s\n", address)
	http.ListenAndServe(address, h)
}
