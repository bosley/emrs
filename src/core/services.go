package core

import (
	"log/slog"
	"sync"
)

type serviceManager struct {
	mu    *sync.Mutex
	known map[string]Service
}

func newServiceManager() *serviceManager {
	return &serviceManager{
		mu:    new(sync.Mutex),
		known: make(map[string]Service),
	}
}

func (s *serviceManager) add(name string, service Service) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.known[name]

	if ok {
		return ErrDuplicateServiceName
	}

	s.known[name] = service
	return nil
}

func (s *serviceManager) start() error {

	s.mu.Lock()
	defer s.mu.Unlock()

	started := make([]Service, 0)

	rollback := func() {
		for _, s := range started {
			if err := s.Stop(); err != nil {
				panic(err.Error())
			}
		}
	}

	for name, service := range s.known {

		slog.Debug("starting service", "name", name)
		if err := service.Start(); err != nil {
			defer rollback()
			return err
		}

		started = append(started, service)
	}

	return nil
}

func (s *serviceManager) stop() error {

	s.mu.Lock()
	defer s.mu.Unlock()

	for name, service := range s.known {
		slog.Debug("stopping service", "name", name)
		if err := service.Start(); err != nil {
			return err
		}
	}
	return nil
}
