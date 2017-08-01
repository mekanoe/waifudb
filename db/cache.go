package db

import (
	"sync"
)

const maxLookup = 1000

type cache struct {
	M          *sync.RWMutex
	Types      map[string]Type
	Q          *sync.RWMutex
	lookup     map[string]map[string]interface{}
	lookupKeys []string
}

func newCache() *cache {
	return &cache{
		M:          &sync.RWMutex{},
		Types:      map[string]Type{},
		Q:          &sync.RWMutex{},
		lookup:     map[string]map[string]interface{}{},
		lookupKeys: make([]string, maxLookup),
	}
}

func (c *cache) GetLookup(ty, k string, v interface{}) (o map[string]interface{}, ok bool) {
	key, err := hashIndex(ty, k, v)
	if err != nil {
		return o, false
	}

	c.Q.RLock()
	o, ok = c.lookup[key]
	c.Q.RUnlock()

	return o, ok
}

func (c *cache) PutLookup(ty, k string, v interface{}, o map[string]interface{}) {
	key, err := hashIndex(ty, k, v)
	if err != nil {
		return
	}

	l := len(c.lookup)

	c.Q.Lock()
	if l > maxLookup {
		delete(c.lookup, c.lookupKeys[0])
		c.lookupKeys = c.lookupKeys[1:]
	}

	c.lookup[key] = o
	c.lookupKeys[l+1] = key
	c.Q.Unlock()
}
