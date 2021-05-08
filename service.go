package main


import (
	"golang.org/x/text/language"
	"time"
)


// User Story/Interactor/App specific business Rules

type RequestModel struct {
	from language.Tag
	to language.Tag
	fromPhrase string
}



// Service is a Translator user.
type Service struct {
	translator            Translator
	cache                 map[RequestModel]string
	// dedublicate service
	busy                  bool
	translatingInProgress setOfRequestModel
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
		cache: make(map[RequestModel]string),
	}
}

