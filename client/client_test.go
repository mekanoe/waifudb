package client

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	run := m.Run()
	os.Exit(run)
}

func TestClientCreate(t *testing.T) {
	_, err := New(nil)
	if err != nil {
		t.Error("create failed", err)
		return
	}
}
