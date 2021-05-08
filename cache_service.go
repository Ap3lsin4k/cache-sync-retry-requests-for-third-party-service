package main

func (s *Service) CacheTranslate(what CacheKey) (string, error) {

	if s.isCached(what) {
		return s.getCachedUnchecked(what)
	}

	return s.translateAndSaveCache(what)
}

/*
func (s *Service) CacheTranslate(what CacheKey) (string, error) {

	value, ok := s.getCachedSafe(what)

	if ok {
		return value
	}

	return s.translateAndSaveCache(what)
}
*/