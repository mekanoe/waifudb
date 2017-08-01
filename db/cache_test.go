package db

import "testing"

func TestCacheLoad(t *testing.T) {
	path := getRandomDBName()
	w, ds := getWaifu(path)

	_, err := w.CreateType(&Type{
		Name:    "testcache",
		Indexes: []string{"name", "instrument"},
	})
	if err != nil {
		t.Error(err)
		return
	}

	tt, err := w.GetType("testcache")
	if err != nil {
		t.Error(err)
		return
	}
	ds.Release()

	w2, ds2 := getWaifu(path)
	tt2, err := w2.GetType("testcache")
	if err != nil {
		t.Error(err)
		return
	}
	ds2.Release()

	if tt.ID != tt2.ID {
		t.Errorf("got two different types on the same key, `%s` != `%s`", tt.Name, tt2.Name)
	}
}
