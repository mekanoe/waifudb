package datastore

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

type testData struct {
	Test1  string
	YesNo  bool
	BiteMe []byte
}

var td = testData{
	Test1:  "なにこれ",
	YesNo:  true,
	BiteMe: []byte("*bites*"),
}

func getDS() (*Datastore, error) {
	return New(&Config{
		Path: "/tmp/.trash." + randStringRunes(10),
	})
}

func freeDS(ds *Datastore) {
	ds.Release()
	os.Remove(ds.Bolt.Path())
}

func TestRWJSONFlow(t *testing.T) {
	ds, err := getDS()
	if err != nil {
		t.Fatal(err)
	}

	err = ds.SetJSON([]byte("data"), "test:1", td)
	if err != nil {
		t.Fatal(err)
	}

	var data testData
	err = ds.GetJSON([]byte("data"), "test:1", &data)
	if err != nil {
		t.Fatal(err)
	}

	freeDS(ds)
}

func TestRWErrors(t *testing.T) {
	ds, err := getDS()
	if err != nil {
		t.Fatal(err)
	}

	err = ds.GetJSON([]byte("NOTGOOD"), "wefj9qwejf", nil)
	if err == nil {
		t.Error("did not error due to bad bucket name")
		return
	}

	d, err := ds.Get([]byte("data"), "notakey")
	if err != ErrNotFound {
		t.Error("notakey not empty, didn't fail. data := %s", d)
		return
	}

	freeDS(ds)
}
