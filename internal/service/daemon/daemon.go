package daemon

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kardianos/service"
)

type Program struct {
	httpHandler http.Handler
	port        int
	stopCh      chan struct{}
}

func (p *Program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *Program) run() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.port),
		Handler: p.httpHandler,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	<-p.stopCh
}

func (p *Program) Stop(s service.Service) error {
	close(p.stopCh)
	return nil
}

type Service struct {
	service service.Service
}

func NewService(name, displayName, description string, httpHandler http.Handler, port int) (*Service, error) {
	svcConfig := &service.Config{
		Name:        name,
		DisplayName: displayName,
		Description: description,
	}

	prg := &Program{
		httpHandler: httpHandler,
		port:        port,
		stopCh:      make(chan struct{}),
	}

	svc, err := service.New(prg, svcConfig)
	if err != nil {
		return nil, fmt.Errorf("create service: %w", err)
	}

	return &Service{
		service: svc,
	}, nil
}

func (s *Service) Run() error {
	return s.service.Run()
}

func (s *Service) Install() error {
	return s.service.Install()
}

func (s *Service) Uninstall() error {
	return s.service.Uninstall()
}

func (s *Service) Start() error {
	return s.service.Start()
}

func (s *Service) Stop() error {
	return s.service.Stop()
}

func (s *Service) Status() (string, error) {
	status, err := s.service.Status()
	if err != nil {
		return "unknown", err
	}

	switch status {
	case service.StatusRunning:
		return "running", nil
	case service.StatusStopped:
		return "stopped", nil
	default:
		return "unknown", nil
	}
}

func WaitForSignal(stopFn func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	if stopFn != nil {
		stopFn()
	}
}
