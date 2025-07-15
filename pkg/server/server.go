package server

import (
	"go_final_project/pkg/api"
	"log"
	"net/http"
	"os"
)

type Server struct {
	logger *log.Logger
	port   string
}

func NewServer(logger *log.Logger) *Server {
	// check the TODO_PORT variable, if empty, then the default server is 7540
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	return &Server{
		logger: logger,
		port:   port,
	}
}

func (s *Server) Start() error {

	// connecting api handlers
	api.Init()
	// take data for the front
	webDir := "./web"

	// create a file server
	fileServer := http.FileServer(http.Dir(webDir))

	// create a handler
	http.Handle("/", fileServer)

	// start the server
	s.logger.Println("The server is running on http://localhost:" + s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
