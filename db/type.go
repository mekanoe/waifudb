package db

import (
	"encoding/json"
	"sync"

	"github.com/segmentio/ksuid"
)

type Type struct {
	ID   string // KSUID
	Name string
}

func (w *WaifuDB) CreateType(name string) error {
	kid, err := ksuid.NewRandom()
	if err != nil {
		return err
	}
	id := kid.String()

	t := Type{
		ID:   id,
		Name: name,
	}

	err = w.store.SetJSON(bktTypes, id, &t)
	if err != nil {
		return err
	}

	w.cache.M.Lock()
	w.cache.Types[name] = t
	w.cache.M.Unlock()

	return nil
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
