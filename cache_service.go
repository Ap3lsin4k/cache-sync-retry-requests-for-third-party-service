package main

import (
	"context"
	"fmt"
)

func (s *Service) CacheTranslate(what RequestModel) (string, error) {

	if s.isCached(what) {
		return s.getCachedUnchecked(what)
	}

	return s.translateAndSaveCache(what)
}

/*
func (s *Service) CacheTranslate(what RequestModel) (string, error) {

	value, ok := s.getCachedSafe(what)

	if ok {
		return value
	}

	return s.translateAndSaveCache(what)
}
*/

func (s Service) isCached(that RequestModel) bool {
	_, ok := s.cache[that]
	return ok
}

func (s *Service) getCachedUnchecked(that RequestModel) (string, error) {
	val, _ := s.cache[that]
	return val, nil
}

func recoverCacheUninitialized() {
	if r := recover(); r != nil {

		//assert r equals "assignment to entry in nil map"
		descriptive := fmt.Errorf(
			"\t try initialize service.cache with `make(map[RequestModel]string)` in the service constructor")

		panic(descriptive)
	}
}

func (s Service) Cache(requestModel RequestModel, responseModel string) {
	s.cache[requestModel] = responseModel
}

func (s *Service) translateAndSaveCache(key RequestModel) (string, error) {
	ctx := context.Background()
	translated, err := s.translator.Translate(ctx, key.from, key.to, key.fromPhrase)
	if err != nil {
		return "", err
	}

	defer recoverCacheUninitialized()
	s.cache[key] = translated

	return translated, nil
}
