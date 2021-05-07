package main

import (
	"context"
	"fmt"
	"golang.org/x/text/language"
	"math/rand"
	"strconv"
	"testing"
)

type translatorCacheSpy struct {
	TranslateCount uint8
}

func (t *translatorCacheSpy) Translate(_ context.Context, from, to language.Tag, data string) (string, error) {
	t.TranslateCount++
	res := fmt.Sprintf("%v -> %v : %v -> %v; service was called %d times", from, to, data, strconv.FormatInt(rand.Int63(), 10), t.TranslateCount)
	return res, nil
}




func NewDummyService() *Service {
	t := &translatorCacheSpy{
		0,
	}


	return &Service{
		translator: t,
	}
}

func TestCacheWhenCallingTwoTimes(t *testing.T){

	s := NewDummyService()

	spy := s.translator.(*translatorCacheSpy)
	if spy.TranslateCount != 0 {
		t.Errorf("Expect `TranslateCount` to be zero by default")
	}
	s.CacheTranslate(language.English, language.Japanese, "same words")

	if spy.TranslateCount != 1 {
		t.Errorf("Expect `Translate` to be called")
	}

	s.CacheTranslate(language.English, language.Japanese, "same words")

	if spy.TranslateCount != 1 {
		t.Errorf("`Translate` must be called once when ServiceTranslate is called twice")
	}
}

func TestCacheWhenCalledWithDifferentInputs(t *testing.T) {

}

// TODO benchmark Golang's native cache and compare with a hash map