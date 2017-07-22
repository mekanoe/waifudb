package db

import (
	"sync"
)

type cache struct {
	M     *sync.RWMutex
	Types map[string]Type
}

func newCache() *cache {
	return &cache{
		M:     &sync.RWMutex{},
		Types: map[string]Type{},
	}
}
