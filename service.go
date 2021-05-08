package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"log"
	"sync"
	"time"
)

type CacheKey struct {
	from language.Tag
	to language.Tag
	fromPhrase string
}

type Promise struct {
	translated string
	err 	error
}

type in struct{}
var active in
// GoLang does have mathematical set implementation so we use map without values (values are empty structs) to emulate that
type setOfRequestModel map[CacheKey]in

// Service is a Translator user.
type Service struct {
	translator            Translator
	cache                 map[CacheKey]string
	promise               Promise
	busy                  bool
	translatingInProgress setOfRequestModel
}

func (s Service) isCached(that CacheKey) bool {
	_, ok := s.cache[that]
	return ok
}

func (s *Service) getCachedUnchecked(that CacheKey) (string, error) {
	val, _ := s.cache[that]
	return val, nil
}

func recoverCacheUninitialized() {
	if r := recover(); r != nil {

		//assert r equals "assignment to entry in nil map"
		descriptive := fmt.Errorf(
			"\t try initialize service.cache with `make(map[CacheKey]string)` in the service constructor")
		

		panic(descriptive)
	}
}

func (s *Service) translateAndSaveCache(key CacheKey) (string, error) {
	ctx := context.Background()
	translated, err := s.translator.Translate(ctx, key.from, key.to, key.fromPhrase)
	if err != nil {
		return "", err
	}

	defer recoverCacheUninitialized()
	s.cache[key] = translated

	return translated, nil
}
// TODO get cache after cache was saved

func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)

	return &Service{
		translator: t,
		cache: make(map[CacheKey]string),
	}
}

func (s *Service) TranslatePromise(ctx context.Context, howTranslate CacheKey) (string, error) {
	//		return s.promise
	// await
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

		/*		v := ctx.Value(howTranslate)
				if v != nil {
					val, ok := ctx.Value(howTranslate).(string)
					if !ok {
						panic("type stored in context must be string")
					}

					return val + "<=", nil
				}
		*/
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

func (s Service) BusyWith(translationRequestID CacheKey) bool {
	_, exists := s.translatingInProgress[translationRequestID]
	return exists
}

func (s Service) AddTranslateInProgress(translate CacheKey, m *sync.Mutex) {
	m.Lock()
	s.translatingInProgress[translate] = active
	m.Unlock()
}

func (s Service) RemoveTranslateProcess(translatingID CacheKey, m *sync.Mutex) {
	m.Lock()
	delete(s.translatingInProgress, translatingID)
	m.Unlock()
}

func (s Service) Cache(requestModel CacheKey, responseModel string) {
	s.cache[requestModel] = responseModel
}