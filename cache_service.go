package main

import (
	"context"
)

func (s *Service) CacheTranslate(what RequestModel) (string, error) {

	value, ok := s.cache.LoadCache(what)

	if ok {
		return value, nil
	}
	ctx := context.Background()

	return s.translateAndSaveCacheAndMightLog(ctx, what)
}


// TODO refactor use Golang like style to work with SaveCache

func (s *Service) tryCache(requestModel RequestModel, responseModel string, err error) error {
	if err != nil {
		return err
	}
	s.cache.SaveCache(requestModel, responseModel)
	return nil
}

