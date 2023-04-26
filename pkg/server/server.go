package server

import (
	"fmt"
	"log"
)

type Server interface {
	Run() error
}

type BehaveHandler func() error

type Servers struct {
	srvs []Server
	befs []BehaveHandler
	afts []BehaveHandler
}

func NewServers() *Servers {
	return &Servers{}
}

func (s *Servers) BindServer(srv Server) {
	s.srvs = append(s.srvs, srv)
}

func (s *Servers) BindBeforeHandler(fn ...BehaveHandler) {
	s.befs = append(s.befs, fn...)
}

func (s *Servers) BindAfterHandler(fn ...BehaveHandler) {
	s.afts = append(s.afts, fn...)
}

func (s *Servers) Run() error {
	if len(s.srvs) == 0 {
		return fmt.Errorf("empty Server")
	}
	for _, fn := range s.befs {
		if err := fn(); err != nil {
			return err
		}
	}

	errChan := make(chan error, len(s.srvs))
	for _, v := range s.srvs {
		go func(s Server) {
			errChan <- s.Run()
		}(v)
	}

	if err := <-errChan; err != nil {
		return log.Default().Output(1, err.Error())
	}

	for _, fn := range s.afts {
		if err := fn(); err != nil {
			return err
		}
	}
	log.Println("server stop")
	return nil
}
