package repositories

import (
	"sync"
	"time"
)

type LoginCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
	Used      bool
}

type MockStorage struct {
	Codes map[string]LoginCode
	Mutex sync.Mutex
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Codes: make(map[string]LoginCode),
	}
}

// Save a new login code
func (s *MockStorage) SaveLoginCode(email, code string, expiresAt time.Time) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Codes[code] = LoginCode{
		Email:     email,
		Code:      code,
		ExpiresAt: expiresAt,
		Used:      false,
	}
}

// Retrieve and delete a login code
func (s *MockStorage) GetLoginCode(code string) (LoginCode, bool) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	loginCode, exists := s.Codes[code]
	if exists {
		delete(s.Codes, code) // Delete the code from the map
	}
	return loginCode, exists
}
