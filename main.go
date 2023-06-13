package main

import (
	"log"
	"os"
	"user_service/startup"
	cfg "user_service/startup/config"
)

func main() {
	log.SetOutput(os.Stderr)
	config := cfg.NewConfig()
	log.Println("Starting server User Service...")
	server := startup.NewServer(config)
	server.Start()
}
