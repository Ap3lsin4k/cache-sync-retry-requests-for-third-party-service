package main

import (
	"time"
)


// User Story/Interactor/App specific business Rules



// Service is a Translator user.
type Service struct {
	translator          Translator // external API
	cache	  			CacheRepository
	// dedublicate service
	dedublicate			DedublicateService
}


func NewService() *Service {
	t := newRandomTranslator(
		100*time.Millisecond,
		500*time.Millisecond,
		0.1,
	)

	return &Service{
		translator: t,
		cache: newMapCacheRepository(),
		dedublicate: newStandardDedublicateService(),
	}
}

