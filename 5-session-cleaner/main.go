//////////////////////////////////////////////////////////////////////
//
// Given is a SessionManager that stores session information in
// memory. The SessionManager itself is working, however, since we
// keep on adding new sessions to the manager our program will
// eventually run out of memory.
//
// Your task is to implement a session cleaner routine that runs
// concurrently in the background and cleans every session that
// hasn't been updated for more than 5 seconds (of course usually
// session times are much longer).
//
// Note that we expect the session to be removed anytime between 5 and
// 7 seconds after the last update. Also, note that you have to be
// very careful in order to prevent race conditions.
//

package main

import (
	"errors"
	"log"
	"sync"
	"time"
)

// SessionManager keeps track of all sessions from creation, updating
// to destroying.
type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]Session
	quit     chan struct{}
}

// Session stores the session's data
type Session struct {
	Data        map[string]interface{}
	LastUpdated time.Time
}

var (
	cleanPeriod = time.Second
	cleanExpiry = 5 * time.Second
)

// NewSessionManager creates a new sessionManager
func NewSessionManager() *SessionManager {
	m := &SessionManager{
		sessions: make(map[string]Session),
		quit:     make(chan struct{}),
	}

	// initialise the session cleaner
	go m.cleanSessions(cleanPeriod, cleanExpiry)

	return m
}

// CreateSession creates a new session and returns the sessionID
func (m *SessionManager) CreateSession() (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sessionID, err := MakeSessionID()
	if err != nil {
		return "", err
	}

	m.sessions[sessionID] = Session{
		Data:        make(map[string]interface{}),
		LastUpdated: time.Now(),
	}

	return sessionID, nil
}

// ErrSessionNotFound returned when sessionID not listed in
// SessionManager
var ErrSessionNotFound = errors.New("SessionID does not exists")

// GetSessionData returns data related to session if sessionID is
// found, errors otherwise
func (m *SessionManager) GetSessionData(sessionID string) (map[string]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, ok := m.sessions[sessionID]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session.Data, nil
}

// UpdateSessionData overwrites the old session data with the new one
func (m *SessionManager) UpdateSessionData(sessionID string, data map[string]interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.sessions[sessionID]
	if !ok {
		return ErrSessionNotFound
	}

	// Hint: you should renew expiry of the session here
	m.sessions[sessionID] = Session{
		Data:        data,
		LastUpdated: time.Now(),
	}

	return nil
}

// Ticker method adapted from: https://stackoverflow.com/a/16466581
func (m *SessionManager) cleanSessions(period time.Duration, expiry time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()

			for key, session := range m.sessions {
				now := time.Now()
				if now.Sub(session.LastUpdated) >= expiry {
					// delete expired session
					delete(m.sessions, key)
				}
			}

			m.mu.Unlock()
		case <-m.quit:
			return
		}
	}
}

func (m *SessionManager) KillSessionCleaner() {
	if m.quit == nil {
		return
	}

	close(m.quit)
	m.quit = nil
}

func main() {
	// Create new sessionManager and new session
	m := NewSessionManager()
	sID, err := m.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created new session with ID", sID)

	// Update session data
	data := make(map[string]interface{})
	data["website"] = "longhoang.de"

	err = m.UpdateSessionData(sID, data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Update session data, set website to longhoang.de")

	// Retrieve data from manager again
	updatedData, err := m.GetSessionData(sID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Get session data:", updatedData)
}
