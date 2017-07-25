package db

import (
	"encoding/json"
	"sync"

	"github.com/segmentio/ksuid"
)

type Type struct {
	ID      string // KSUID
	Name    string
	Indexes []string
}

func (w *WaifuDB) CreateType(name string, indexes []string) (*Type, error) {
	kid, err := ksuid.NewRandom()
	if err != nil {
		return nil, err
	}
	id := kid.String()

	t := Type{
		ID:      id,
		Name:    name,
		Indexes: indexes,
	}

	err = w.store.SetJSON(bktTypes, id, &t)
	if err != nil {
		return nil, err
	}

	w.cache.M.Lock()
	w.cache.Types[name] = t
	w.cache.M.Unlock()

	return &t, nil
}

func (w *WaifuDB) GetType(name string) (*Type, error) {

	w.cache.M.RLock()
	t, ok := w.cache.Types[name]
	w.cache.M.RUnlock()
	if !ok {
		return nil, ErrTypeNotFound
	}

	return &t, nil
}

func (w *WaifuDB) loadTypes() error {
	w.cache.M.Lock()
	defer w.cache.M.Unlock()
	ct := w.cache.Types

	m := &sync.Mutex{}
	return w.store.Walk(bktTypes, func(k, v []byte) error {
		m.Lock()
		defer m.Unlock()

		var data Type
		json.Unmarshal(v, &data)

		ct[data.Name] = data

		return nil
	})
}

func (t *Type) HasIndex(i string) bool {
	for _, v := range t.Indexes {
		if i == v {
			return true
		}
	}
	return false
}
