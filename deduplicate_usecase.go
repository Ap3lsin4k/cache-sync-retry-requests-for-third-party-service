package main

import (
	"context"
	"log"
	"sync"
	"time"
)


type in struct{}
var active in
// GoLang does have mathematical set implementation so we use map without values (values are empty structs) to emulate that
type setOfRequestModel map[RequestModel]in


func (s *Service) TranslatePromise(ctx context.Context, howTranslate RequestModel) (string, error) {
	log.Println("context", ctx.Value("whereami"), "BusyWith?", s.BusyWith(howTranslate))
	if s.BusyWith(howTranslate) {
		for s.BusyWith(howTranslate) {
			time.Sleep(100 * time.Millisecond)
		}
		if !s.isCached(howTranslate) {
			panic("WTF?! cache must be set after resource was removed")
		}
		// FIXME after it was busy it might return with error
		return s.getCachedUnchecked(howTranslate)
	} else {
		var translProcessesMutex = &sync.Mutex{}

		s.AddTranslateInProgress(howTranslate, translProcessesMutex)
		log.Println("context", ctx.Value("whereami"), "\tAddTranslateInProgress", s.BusyWith(howTranslate))

		ctx_translation_in_progress := context.WithValue(ctx, howTranslate, "")

		transl, err := s.translator.Translate(ctx_translation_in_progress, howTranslate.from, howTranslate.to, howTranslate.fromPhrase)
		log.Println("context", ctx.Value("whereami"), "\ttranslator.Translate has been finished")

		s.Cache(howTranslate, transl)
		log.Printf("context: %v\tresult \"%s\" has been cached: %v",
			ctx.Value("whereami"), transl, s.isCached(howTranslate))
		s.RemoveTranslateProcess(howTranslate, translProcessesMutex)
		log.Println("context", ctx.Value("whereami"), "removed translate process, BusyWith:", s.BusyWith(howTranslate))

		return transl, err
	}
}

func (s Service) BusyWith(translationRequestID RequestModel) bool {
	_, exists := s.translatingInProgress[translationRequestID]
	return exists
}

func (s Service) AddTranslateInProgress(translate RequestModel, m *sync.Mutex) {
	m.Lock()
	s.translatingInProgress[translate] = active
	m.Unlock()
}

func (s Service) RemoveTranslateProcess(translatingID RequestModel, m *sync.Mutex) {
	m.Lock()
	delete(s.translatingInProgress, translatingID)
	m.Unlock()
}

