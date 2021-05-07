package main

import (
	"context"
	"golang.org/x/text/language"
)

func (s Service) CacheTranslate(from, to language.Tag, phrase string) {

	ctx := context.Background()
	if isCached(from, to, phrase)
		return GetCachedUnchecked(from, to, phrase)
	return s.translator.Translate(ctx, language.English, language.Japanese, "test")
}

