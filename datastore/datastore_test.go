package datastore

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	r := m.Run()
	os.Remove(".trash.db")
	os.Remove(".trash2.db")
	os.Exit(r)
}

func TestCreateDatastore(t *testing.T) {
	_, err := New(nil)
	if err != nil {
		t.Errorf("failed to create datastore: %v", err)
		t.Fail()
	}

	_, err = New(&Config{
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
}
