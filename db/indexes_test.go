package db

import "testing"

func TestIndexGet(t *testing.T) {
	w, ds := getWaifu("")

	_, err := w.CreateType("person", []string{"name"})
	if err != nil {
		t.Error(err)
		return
	}

	ds.Release()
}

func BenchmarkIndexes(b *testing.B) {
	w, ds := getWaifu(".benchmark-1.db")

	ty, err := w.CreateType("name", []string{"name", "instrument"})
	if err != nil {
		b.Error(err)
		return
	}

	seed(w, ty, 10000, 5)

	knownData := map[string]interface{}{
		"name":       "Reina Kousaka",
		"instrument": "trumpet",
		"loves":      "me",
	}

	w.PutItem(ty.Name, knownData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.GetItemByKey(ty.Name, "name", knownData["name"])
		if err != nil {
			b.Error(err)
			return
		}
	}

	ds.Release()
	ds.DestroyDestroyDestroy()
}

func BenchmarkNoIndexes(b *testing.B) {
	w, ds := getWaifu(".benchmark-2.db")

	ty, err := w.CreateType("name", []string{"name", "instrument"})
	if err != nil {
		b.Error(err)
		return
	}

	knownData := map[string]interface{}{
		"name":       "Reina Kousaka",
		"instrument": "trumpet",
		"loves":      "me",
	}

	w.PutItem(ty.Name, knownData)

	seed(w, ty, 10, 5)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := w.GetItemByKey(ty.Name, "loves", knownData["loves"])
		if err != nil {
			b.Error(err)
			return
		}
	}

	ds.Release()
	ds.DestroyDestroyDestroy()
}
