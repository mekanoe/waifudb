package db

import (
	"errors"

	"github.com/kayteh/waifudb/datastore"
)

type WaifuDB struct {
	store *datastore.Datastore

	cache *cache
}

var (
	ErrTypeNotFound = errors.New("waifudb: type not found")
)

var (
	bktData     = []byte("data")
	bktInternal = []byte("internal")
	bktTypes    = []byte("types")
	bktIndexes  = []byte("indexes")
)

func New(store *datastore.Datastore) (*WaifuDB, error) {
	w := &WaifuDB{
		store: store,
		cache: newCache(),
	}

	err := w.loadTypes()
	if err != nil {
		return nil, err
	}

	return w, nil
}
