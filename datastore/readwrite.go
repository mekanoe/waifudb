package datastore

import (
	"bytes"
	"encoding/json"

	"github.com/boltdb/bolt"
)

var (
	managedBuckets = [][]byte{
		[]byte("data"),
		[]byte("types"),
		[]byte("internal"),
		[]byte("indexes"),
	}
)

func checkBucket(name []byte) error {
	good := false

	for _, n := range managedBuckets {
		if good {
			break
		}

		if bytes.Compare(n, name) == 0 {
			good = true
			break
		}
	}

	if good {
		return nil
	}

	return ErrBadBucket
}

// Set some data
func (ds *Datastore) Set(bucket []byte, key string, val []byte) error {
	err := checkBucket(bucket)
	if err != nil {
		return err
	}

	return ds.Bolt.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Put([]byte(key), val)
	})
}

// SetJSON some json data in a bucket
func (ds *Datastore) SetJSON(bucket []byte, key string, val interface{}) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	return ds.Set(bucket, key, data)
}

// Get some data
func (ds *Datastore) Get(bucket []byte, key string) ([]byte, error) {
	err := checkBucket(bucket)
	if err != nil {
		return nil, err
	}

	var out []byte

	err = ds.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		out = b.Get([]byte(key))
		if len(out) == 0 {
			return ErrNotFound
		}
		return nil
	})
	return out, err
}

// GetJSON auto-unmarshals data from a bucket's key
func (ds *Datastore) GetJSON(bucket []byte, key string, out interface{}) error {
	d, err := ds.Get(bucket, key)
	if err != nil {
		return err
	}

	return json.Unmarshal(d, out)
}

// Walk through the bucket, running `fn` every iteration
func (ds *Datastore) Walk(bucket []byte, fn func([]byte, []byte) error) error {
	err := checkBucket(bucket)
	if err != nil {
		return err
	}

	return ds.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.ForEach(fn)
	})
}
