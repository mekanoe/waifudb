package db

import (
	"fmt"
	"testing"

	"github.com/boltdb/bolt"
)

func TestQuads(t *testing.T) {
	w, ds := getWaifu("")

	_, err := w.CreateType(&Type{
		Name:    "quads",
		Indexes: []string{"name"},
		Relations: map[string]string{
			"friend": "friend",
			"asym1":  "asym2",
			"enemy":  "",
		},
	})
	if err != nil {
		t.Error("create type failed", err)
		return
	}

	test := map[string]interface{}{
		"name":       "Reina Kousaka",
		"instrument": "Trumpet",
		"friend":     "quads:0rdXHIt541pf5DCazwCk3rEHffK",
		"asym1":      "quads:0rdXJJPGd3BFsLiDOhgTpYwxrJs",
		"asym2":      "quads:0rdXKwnwLbFNM9bxkvZMRATOCBS",
		"enemy":      "quads:0rdXMj4BGHPPwPevsLr7lfNnzc4",
	}

	i, err := w.PutItem("quads", test)
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(i)

	w.store.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bktQuads)

		b.ForEach(func(k []byte, v []byte) error {
			fmt.Printf("%s => %s\n", k, v)
			return nil
		})

		return nil
	})

	ds.Release()
}
