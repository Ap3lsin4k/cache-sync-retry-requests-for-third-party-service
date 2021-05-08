package main

import (
	"context"
	"log"
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
		log.Println("context", ctx.Value("whereami"), "_BusyWith?", s.BusyWith(howTranslate))

		val, ok := s.getCachedSafe(howTranslate)
		log.Println("context", ctx.Value("whereami"), "getCached?", ok)
		if ok {
			return val, nil
		}
		// if cache is not set than we just try our own attempts
	}

	s.AddTranslateInProgress(howTranslate)
	defer s.RemoveTranslateProcess(howTranslate)

	log.Println("context", ctx.Value("whereami"), "\tAddTranslateInProgress", s.BusyWith(howTranslate))


	transl, err := s.translator.Translate(ctx, howTranslate.from, howTranslate.to, howTranslate.fromPhrase)
	log.Println("context", ctx.Value("whereami"), "\ttranslator.Translate has been finished with (res:",transl,", err:",err,")")

	if err != nil {
		return "", err
	}
	s.Cache(howTranslate, transl)
	log.Printf("context: %v\tresult \"%s\" has been cached: %v",
		ctx.Value("whereami"), transl, s.isCached(howTranslate))


	log.Println("context", ctx.Value("whereami"), "removed translate process, BusyWith:", s.BusyWith(howTranslate))

	return transl, err

}

func (s *Service) BusyWith(translationRequestID RequestModel) bool {
	s.mutex.Lock()
	_, exists := s.translatingInProgress[translationRequestID]
	s.mutex.Unlock()
	return exists
}

func (s *Service) AddTranslateInProgress(translate RequestModel) {
	s.mutex.Lock()
	s.translatingInProgress[translate] = active
	s.mutex.Unlock()
}

func (s *Service) RemoveTranslateProcess(translatingID RequestModel) {
	s.mutex.Lock()
	delete(s.translatingInProgress, translatingID)
	s.mutex.Unlock()
}

