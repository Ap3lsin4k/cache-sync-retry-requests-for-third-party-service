package main

import (
	"sync"
	"time"
)

func (s *StandardDedublicateService) waitForResourceToBeAvailable(howTranslate RequestModel) {
	for s.BusyWith(howTranslate) {
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *StandardDedublicateService) BusyWith(translationRequestID RequestModel) bool {
	s.mutex.Lock()
	_, exists := s.translatingInProgress[translationRequestID]
	s.mutex.Unlock()
	return exists
}

func (s *StandardDedublicateService) AddTranslateInProgress(translate RequestModel) {
	s.mutex.Lock()
	s.translatingInProgress[translate] = active
	s.mutex.Unlock()
}

func (s *StandardDedublicateService) RemoveTranslateProcess(translatingID RequestModel) {
	s.mutex.Lock()
	delete(s.translatingInProgress, translatingID)
	s.mutex.Unlock()
}

type StandardDedublicateService struct {
	busy                  bool
	translatingInProgress setOfRequestModel
	mutex                 sync.Mutex
}

func newStandardDedublicateService() *StandardDedublicateService {
	return &StandardDedublicateService{
		busy:                  false,
		translatingInProgress: make(setOfRequestModel),
		mutex:                 sync.Mutex{},
	}
}