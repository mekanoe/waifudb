package datastore

import (
	"os"
	"sync"
	"testing"

	"bytes"

	"github.com/boltdb/bolt"
)

func TestMain(m *testing.M) {
	r := m.Run()
	os.Remove(".trash.db")
	os.Remove(".trash2.db")
	os.Exit(r)
}

func TestCreateDatastore(t *testing.T) {
	ds, err := getDS()
	if err != nil {
		t.Errorf("failed to create datastore: %v", err)
		t.Fail()
	}
	freeDS(ds)

	ds, err = New(&Config{
		Path: ".trash2.db",
	})
	if err != nil {
		t.Errorf("failed to create datastore: %v", err)
		t.Fail()
	}

	_, err = os.Stat(".trash2.db")
	if os.IsNotExist(err) {
		t.Errorf(".trash2.db doesn't exist, config failed.")
		t.Fail()
	}

	freeDS(ds)
}

func seedBolt(data map[string]string, db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("data"))
		for k, v := range data {
			err := b.Put([]byte(k), []byte(v))
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func TestWalk(t *testing.T) {
	store, err := getDS()
	if err != nil {
		t.Error(err)
		return
	}
	seed := map[string]string{
		"abc": "false",
		"def": "false",
		"ghi": "false",
	}

	seedBolt(seed, store.Bolt)

	m := &sync.Mutex{}
	check := [][]byte{}

	err = store.Walk([]byte("data"), func(k, v []byte) error {
		m.Lock()
		check = append(check, k)
		m.Unlock()
		return nil
	})

	if err != nil {
		t.Error(err)
		return
	}

	checkStr := string(bytes.Join(check, []byte("")))
	if checkStr != "abcdefghi" {
		t.Errorf("expected `abcdefghi` got `%s`", checkStr)
		return
	}
}
