package run

import (
	"net/http"
	"testing"

	"encoding/json"

	"github.com/hydrogen18/memlistener"
	"github.com/valyala/fasthttp"
)

func getInmemoryServer() (*memlistener.MemoryListener, *http.Client) {
	iml := memlistener.NewMemoryListener()

	s := &fasthttp.Server{
		Handler: handler,
	}

	tport := &http.Transport{}
	tport.Dial = iml.Dial

	client := &http.Client{
		// Force no redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},

		Transport: tport,
	}

	go s.Serve(iml)

	return iml, client
}

func TestStart(t *testing.T) {
	go Start("unix:/tmp/test-waifu")
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

func TestEmptyQuery(t *testing.T) {
	iml, cl := getInmemoryServer()

	rq, err := http.NewRequest("QUERY", "http://localhost/", nil)
	if err != nil {
		t.Error("request create error: ", err)
		return
	}

	_, err = cl.Do(rq)
	if err != nil {
		t.Error("request create error: ", err)
		return
	}

	// var resp OutgoingMessage
	// err = json.NewDecoder(rsp.Body).Decode(&resp)
	// if err != nil {
	// 	t.Error("decoder error: ", err)
	// 	return
	// }

	// if !resp.Success || resp.Payload.(string) != "resp" {
	// 	t.Errorf("resp response is bad: %v", resp)
	// }

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

	if msg.Success || msg.Payload.(string) != "bad request" {
		t.Errorf("msg response is bad: %v", msg)
	}

	iml.Close()
}
