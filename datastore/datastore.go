package datastore

import (
	"errors"
	"log"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/kayteh/waifudb/util"
)

var (
	ErrGeneralError  = errors.New("waifudb/datastore: general error")
	ErrSettingsError = errors.New("waifudb/datastore: settings error")
)

type Datastore struct {
	lock *sync.RWMutex
	bolt *bolt.DB
}

func NewDatastore() (*Datastore, error) {
	path, err := util.Getenvdef("DATA_PATH", ".trash/bolt.db").String()
	if err != nil {
		log.Fatalln("failed to get DATA_PATH", err)
		return nil, ErrSettingsError
	}

	b, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatalf("failed to open %s, %v\n", path, err)
		return nil, err
	}

	var lock *sync.RWMutex
	ds := &Datastore{
		lock: lock,
		bolt: b,
	}

	return ds, nil
}
