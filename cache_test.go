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
		cache: make(map[RequestModel]string),
	}
}

func TestCacheWhenCallingTwoTimes(t *testing.T){

	s := NewDummyService()

	spy := s.translator.(*translatorCacheSpy)
	if spy.TranslateCount != 0 {
		t.Errorf("Expect `TranslateCount` to be zero by default")
	}

	useCase1 := RequestModel{
		from:       language.English,
		to:         language.Japanese,
		fromPhrase: "same words",
	}


	_, _ = s.CacheTranslate(useCase1)
	if !s.isCached(useCase1) {
		t.Errorf("Must cache %v", useCase1)
	}

	if spy.TranslateCount != 1 {
		t.Errorf("Expect `Translate` to be called")
	}

	_, _ = s.CacheTranslate(useCase1)

	if spy.TranslateCount != 1 {
		t.Errorf("`Translate` must be called once when ServiceTranslate is called twice")
	}
}

func TestNotUsingCacheWhenInputDiffers(t *testing.T) {
	s := NewDummyService()
	remoteTranslateCount := &(s.translator.(*translatorCacheSpy)).TranslateCount

	if *remoteTranslateCount != 0 {
		t.Errorf("Expect `TranslateCount` to be zero by default")
	}

	useCase1 := RequestModel{
		from:       language.English,
		to:         language.Armenian,
		fromPhrase: "apple",
	}

	useCase2 := RequestModel{
		from:       language.English,
		to:         language.Armenian,
		fromPhrase: "banana",
	}
	useCase3 := RequestModel{
		from:       language.English,
		to:         language.Polish,
		fromPhrase: "horse",
	}

	_, _ = s.CacheTranslate(useCase1)

	if *remoteTranslateCount != 1 {
		t.Errorf("Expect `TranslateCount` to be 1, input:%v", useCase1)
	}
	_, _ = s.CacheTranslate(useCase2)

	if *remoteTranslateCount != 2 {
		t.Errorf("Expect 'TranslateCount' to be 2, but got %d, when called Translate(%v)",
			*remoteTranslateCount, useCase2)
	}

	_, _ = s.CacheTranslate(useCase3)

	if *remoteTranslateCount != 3 {
		t.Errorf("Expect 'TranslateCount' to be 3, but got %d, when called Translate(%v), i.e. same `fromPhrase` but different language ",
			*remoteTranslateCount, useCase3)
	}
}

// TODO use Context for cache
// TODO benchmark Golang's native cache and compare with a hash map