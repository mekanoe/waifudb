package datastore

import (
	"errors"
	"log"
	"sync"

	"os"

	"github.com/boltdb/bolt"
	"github.com/imdario/mergo"
)

var (
	// ErrGeneralError is a generic error if we can't be more specific
	ErrGeneralError = errors.New("waifudb/datastore: general error")

	// ErrSettingsError happens on any failure related to settings
	ErrSettingsError = errors.New("waifudb/datastore: settings error")

	// ErrBadBucket means you used a bucket that this datastore doesn't deal with.
	ErrBadBucket = errors.New("waifudb/datastore: bucket not relevant")

	// ErrNotFound occurs on unset/empty keys.
	ErrNotFound = errors.New("waifudb/datastore: no value at requested key")

	// ErrStopWalking ends a walk
	ErrStopWalking = errors.New("walk stop")
)

// Datastore contains a RWMutex and a Bolt.DB instance for all the
// sotrage goodness we'll ever need
type Datastore struct {
	lock *sync.RWMutex
	Bolt *bolt.DB
}

// Config is an optional override struct
type Config struct {
	// Path to database file. The folder ideally should exist.
	Path string

	// FileMode is the filemode of the DB file, e.g. 0600 is valid.
	FileMode os.FileMode
}

func (c *Config) merge(incoming *Config) error {
	if incoming == nil {
		return nil
	}

	return mergo.MergeWithOverwrite(c, incoming)
}

var (
	defaultConfig = &Config{
		Path:     ".trash.db",
		FileMode: 0600,
	}
)

// New creates a Datastore
func New(cfg *Config) (*Datastore, error) {
	c := defaultConfig
	err := c.merge(cfg)
	if err != nil {
		return nil, err
	}

	b, err := bolt.Open(c.Path, c.FileMode, nil)
	if err != nil {
		log.Fatalf("failed to open %s, %v\n", c.Path, err)
		return nil, err
	}

	var lock *sync.RWMutex
	ds := &Datastore{
		lock: lock,
		Bolt: b,
	}

	err = b.Update(func(tx *bolt.Tx) error {

		for _, bucket := range managedBuckets {
			_, err := tx.CreateBucketIfNotExists(bucket)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("failed to create buckets: %v", err)
	}

	return ds, nil
}

func (ds *Datastore) Release() error {
	return ds.Bolt.Close()
}

func (ds *Datastore) DestroyDestroyDestroy() error {
	return os.Remove(ds.Bolt.Path())
}
