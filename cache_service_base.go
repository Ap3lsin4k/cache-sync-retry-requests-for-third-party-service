package main

import (
	"context"
	"fmt"
	"log"
)

func recoverCacheUninitialized() {
	if r := recover(); r != nil {

		//assert r equals "assignment to entry in nil map"
		descriptive := fmt.Errorf(
			"\t try initialize service.cache with `make(map[RequestModel]string)` in the service constructor")

		panic(descriptive)
	}
}

func (s *Service) translateAndSaveCacheAndMightLog(ctx context.Context, howTranslate RequestModel) (string, error) {
	transl, err := s.translator.Translate(ctx, howTranslate.from, howTranslate.to, howTranslate.fromPhrase)
	log.Println("context", ctx.Value("whereami"), "\ttranslator.Translate has been finished with (res:", transl, ", err:", err, ")")
	defer recoverCacheUninitialized()
	err = s.tryCache(howTranslate, transl, err)
	log.Printf("context: %v\tresult \"%s\" has been cached: %v",
		ctx.Value("whereami"), transl, s.cache.IsCached(howTranslate))

	return transl, err

}