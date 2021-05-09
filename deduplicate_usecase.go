package main

import (
	"context"
	"log"
)


type in struct{}
var active in
// GoLang does have mathematical set implementation so we use map without values (values are empty structs) to emulate that
type setOfRequestModel map[RequestModel]in


func (s *Service) TryTranslate(ctx context.Context, howTranslate RequestModel) (string, error) {
	log.Println("context", ctx.Value("whereami"), "BusyWith?", s.dedublicate.BusyWith(howTranslate))
	return s.TranslateOnce(ctx, howTranslate)
}

func (s *Service) TranslateOnce(ctx context.Context, howTranslate RequestModel) (string, error) {
	if s.dedublicate.BusyWith(howTranslate) {
		s.dedublicate.waitForResourceToBeAvailable(howTranslate)
		log.Println("context", ctx.Value("whereami"), "_BusyWith?", s.dedublicate.BusyWith(howTranslate))

		val, ok := s.cache.LoadCache(howTranslate)
		log.Println("context", ctx.Value("whereami"), "getCached?", ok)
		if ok {
			return val, nil
		}
		// if cache is not set than we just try our own attempts
	}

	s.dedublicate.AddTranslateInProgress(howTranslate)
	defer s.dedublicate.RemoveTranslateProcess(howTranslate)

	log.Println("context", ctx.Value("whereami"), "\tAddTranslateInProgress", s.dedublicate.BusyWith(howTranslate))

	transl, err := s.mytranslateAndSaveCache(ctx, howTranslate)
	return transl, err

}

type DedublicateService interface {
	waitForResourceToBeAvailable(RequestModel)
	BusyWith(RequestModel) bool
	AddTranslateInProgress(RequestModel)
	RemoveTranslateProcess(RequestModel)
}



