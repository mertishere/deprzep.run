package main

import (
	"os"

	server "XerWoho/zep.run/internal/app/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if err := server.RunServer(); err != nil {
		os.Exit(1)
		return
	}
}
