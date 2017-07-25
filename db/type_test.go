package db

import "testing"

func TestTypes(t *testing.T) {
	w, ds := getWaifu("")

	_, err := w.CreateType("test-type", []string{})
	if err != nil {
		t.Error("create type failed", err)
		return
	}

	ty, err := w.GetType("test-type")
	if err != nil {
		t.Error("test-type didn't get cached!", err)
		return
	}

	if ty.Name != "test-type" {
		t.Errorf("what? type `test-type` was actually `%s`", ty.Name)
		return
	}

	ds.Release()
}
