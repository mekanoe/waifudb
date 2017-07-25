package db

import (
	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/kayteh/waifudb/datastore"
)

type WaifuDB struct {
	store  *datastore.Datastore
	logger *logrus.Entry
	cache  *cache
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
		store:  store,
		cache:  newCache(),
		logger: logrus.WithFields(logrus.Fields{}),
	}

	err := w.loadTypes()
	if err != nil {
		return nil, err
	}

	return w, nil
}
