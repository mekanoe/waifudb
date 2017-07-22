package db

import (
	"log"
	"os"
	"testing"

	"github.com/kayteh/waifudb/datastore"
)

func getWaifu() (*WaifuDB, *datastore.Datastore) {
	store, err := datastore.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	w, err := New(store)
	if err != nil {
		log.Fatal(w)
	}

	return w, store
}

func TestMain(m *testing.M) {
	r := m.Run()
	os.Remove(".trash.db")
	os.Remove(".trash2.db")
	os.Exit(r)
}

func TestBasic(t *testing.T) {
	store, err := datastore.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	w, err := New(store)
	if err != nil {
		t.Fatal(w)
	}

	store.Release()
}
