package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
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

// Service is a Translator user.
type Service struct {
	translator Translator
	cache      map[CacheKey]string
	promise    Promise
	busy       bool
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
	s.busy = true

	v := ctx.Value(howTranslate)
	if v != nil {
		val, ok := ctx.Value(howTranslate).(string)
		if !ok {
			panic("type stored in context must be string")
		}

		return val+"<=", nil
	}


	ctx_translation_in_progress := context.WithValue(ctx, howTranslate, "")



/*	if s.busy {
		if s.busyWith(howTranslate) {
			for s.busy {
				time.Sleep(500 * time.Millisecond)
			}
			return s.getCachedUnchecked(howTranslate) // "*", nil // await
		}
	}
	else {
		// busy with different or not busy at all
//		return s.translateAndSaveCache()
	}*/
	transl, err := s.translator.Translate(ctx_translation_in_progress, howTranslate.from, howTranslate.to, howTranslate.fromPhrase)

	s.busy = false
//	ctx_translation := context.WithValue(ctx_translation_in_progress, howTranslate, transl)

	s.promise = Promise{
		translated: transl,
		err:        err,
	}
	return transl, nil
}

func (s Service) BusyWith(CacheKey) bool {
	return s.busy
}