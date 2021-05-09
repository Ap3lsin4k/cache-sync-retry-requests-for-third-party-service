package main

type mapCacheRepository struct {
	cache map[RequestModel]string
}


func (s mapCacheRepository) IsCached(that RequestModel) bool {
	_, ok := s.cache[that]
	return ok
}

func (s mapCacheRepository) LoadCache(what RequestModel) (string, bool) {
	val, ok := s.cache[what]
	return val, ok
}

func (s *mapCacheRepository) SaveCache(requestModel RequestModel, responseModel string) {
	s.cache[requestModel] = responseModel
}

func newMapCacheRepository() *mapCacheRepository {
	return &mapCacheRepository{
		make(map[RequestModel]string),
	}
}

