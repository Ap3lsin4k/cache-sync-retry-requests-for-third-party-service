package main

import (
	"context"
	"log"
)

type CacheRepository interface {
	IsCached(that RequestModel) bool
	LoadCache(what RequestModel) (string, bool)
	SaveCache(requestModel RequestModel, responseModel string)
}

func (s *Service) mytranslateAndSaveCache(ctx context.Context, howTranslate RequestModel) (string, error) {
	transl, err := s.translateAndSaveCacheAndMightLog(ctx, howTranslate)

	if err != nil {
		return "", err
	}
	s.cache.SaveCache(howTranslate, transl)

	log.Println("context", ctx.Value("whereami"), "removed translate process, BusyWith:", s.dedublicate.BusyWith(howTranslate))
	return transl, err
}
