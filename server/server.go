package server

import (
	"event_ticket_booking/config"
	commonModel "event_ticket_booking/model"
	"fmt"
)

type Server struct {
	Cfg config.Config
	Lib commonModel.Lib
}

func NewServer(config config.Config, lib commonModel.Lib) Server {
	return Server{
		Cfg: config,
		Lib: lib,
	}
}

func (s *Server) Start() error {
	router := NewRouter(s.Cfg, s.Lib)
	fmt.Println("Starting server on port", s.Cfg.Server.Port)
	domain := fmt.Sprintf("localhost:%s", s.Cfg.Server.Port)
	return router.Run(domain)
}

func (s *Server) Stop() {
	fmt.Println("Server stopped")
}
