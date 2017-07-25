package run

import (
	"log"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"encoding/json"

	"github.com/hydrogen18/memlistener"
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

func getInmemoryServer() (*memlistener.MemoryListener, *http.Client) {
	iml := memlistener.NewMemoryListener()

	s, err := New(&Config{
		DBPath: getRandomDBName(),
	})
	if err != nil {
		log.Fatalln("server failed to create")
	}

	go s.Serve(iml)

	tport := &http.Transport{}
	tport.Dial = iml.Dial

	client := &http.Client{
		// Force no redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},

		Transport: tport,
	}

	return iml, client
}

func TestStart(t *testing.T) {
	s, err := New(&Config{
		DBPath: getRandomDBName(),
	})
	if err != nil {
		t.Error(err)
		return
	}

	go s.ListenAndServe()
}

func TestPing(t *testing.T) {
	iml, cl := getInmemoryServer()

	rq, err := http.NewRequest("PING", "http://localhost/", nil)
	if err != nil {
		t.Error("request create error: ", err)
		return
	}

	rsp, err := cl.Do(rq)
	if err != nil {
		t.Error("request create error: ", err)
		return
	}

	var pong OutgoingMessage
	err = json.NewDecoder(rsp.Body).Decode(&pong)
	if err != nil {
		t.Error("decoder error: ", err)
		return
	}

	if !pong.Success || pong.Payload.(string) != "pong" {
		t.Errorf("PONG response is bad: %v", pong)
	}

	iml.Close()
}

func TestBadRequest(t *testing.T) {
	iml, cl := getInmemoryServer()

	rq, err := http.NewRequest("GET", "http://localhost/", nil)
	if err != nil {
		t.Error("request create error: ", err)
		return
	}

	rsp, err := cl.Do(rq)
	if err != nil {
		t.Error("request create error: ", err)
		return
	}

	var msg OutgoingMessage
	err = json.NewDecoder(rsp.Body).Decode(&msg)
	if err != nil {
		t.Error("decoder error: ", err)
		return
	}

	if msg.Success || msg.Error != "bad request" {
		t.Errorf("msg response is bad: %v", msg)
	}

	iml.Close()
}
