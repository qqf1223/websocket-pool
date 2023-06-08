package server

import (
	"log"
	"net/http"
	"time"
)

type WebServer struct {
	Handler http.Handler
	Name    string
	Addr    string
	Timeout int
}

func (s *WebServer) Run() error {
	srv := &http.Server{
		Addr:         s.Addr,
		Handler:      s.Handler,
		ReadTimeout:  time.Duration(s.Timeout) * time.Second,
		WriteTimeout: time.Duration(s.Timeout) * time.Second,
	}

	log.Printf("[INFO] http server start, listen: %s\n", s.Addr)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
