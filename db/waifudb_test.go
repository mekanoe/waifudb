package db

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"fmt"

	"github.com/kayteh/waifudb/datastore"
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

func getRandomDBName() string {
	return "/tmp/" + randStringRunes(5) + ".test-waifu.db"
}

func getWaifu(path string) (*WaifuDB, *datastore.Datastore) {
	if path == "" {
		path = getRandomDBName()
	}

	store, err := datastore.New(&datastore.Config{
		Path: path,
	})
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

func TestSeed(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}

	w, db := getWaifu("")

	ty, err := w.CreateType("test", []string{"nani"})
	if err != nil {
		t.Error(err)
		return
	}

	seed(w, ty, 100, 1)

	db.Release()
	db.DestroyDestroyDestroy()
}

func seed(w *WaifuDB, ty *Type, entries int, nonIndexKeys int) {
	for i := 0; i < entries; i++ {
		// if i%100 == 0 {
		// 	fmt.Printf("seeding: %d -> %d\n", entries, i)
		// }
		builtMap := map[string]interface{}{}
		for _, v := range ty.Indexes {
			builtMap[v] = randStringRunes(rand.Intn(64))
		}

		for k := 0; k < nonIndexKeys; k++ {
			key := fmt.Sprintf("nik_%d", k)
			builtMap[key] = randStringRunes(rand.Intn(64))
		}

		_, err := w.PutItem(ty.Name, builtMap)
		if err != nil {
			w.logger.WithError(err).Error("seed failed")
			return
		}
	}
}
