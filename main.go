package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
	"os"
)

func main() {
	// create a logger
	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags)

	// connect to the database
	if err := db.Init("scheduler.db", logger); err != nil {
		logger.Fatal("error initializing database:", err)
	}

	// create and configure the server
	srv := server.NewServer(logger)

	// start the server
	if err := srv.Start(); err != nil {
		logger.Fatal("Error starting server: ", err)
	}
}
