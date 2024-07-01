package core

import (
	"errors"
	"fmt"
	"sync"
)

type SessionData struct {
	Identity      string // badger identity
	AppActivePath string // application path to show for session
}

type sessionManager struct {
	active map[string]SessionData
	mu     *sync.Mutex
}

// TODO: Consider "save session" and "load saved session"

func newSessionManager() *sessionManager {
	return &sessionManager{
		active: make(map[string]SessionData),
		mu:     new(sync.Mutex),
	}
}

func (mgr *sessionManager) newSession(name string, data SessionData) error {

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	_, exists := mgr.active[name]
	if exists {
		return errors.New(fmt.Sprintf("duplicate session name: %s", name))
	}

	mgr.active[name] = data
	return nil
}

func (mgr *sessionManager) deleteSession(name string) {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	delete(mgr.active, name)
}

func (mgr *sessionManager) getSessionData(name string) *SessionData {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	data, exists := mgr.active[name]
	if exists {
		return &data
	}
	return nil
}

func (mgr *sessionManager) updateSession(name string, data SessionData) error {

	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	_, exists := mgr.active[name]
	if exists {
		mgr.active[name] = data
		return nil
	}
	return errors.New(fmt.Sprintf("invalid session name: %s", name))
}
