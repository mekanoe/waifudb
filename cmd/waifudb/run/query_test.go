package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

func TestGetSet(t *testing.T) {
	iml, cl := getInmemoryServer()

	testDataSet := map[string]map[string]interface{}{
		"payload": map[string]interface{}{
			"@a":   "set",
			"type": "person",
			"data": map[string]interface{}{
				"name":       "Reina Kousaka",
				"instrument": "trumpet",
				"loves":      "me",
			},
		},
	}

	testDataGet := map[string]map[string]interface{}{
		"payload": map[string]interface{}{
			"@a":   "get",
			"type": "person",
			"id":   "",
		},
	}

	testDataType := map[string]map[string]interface{}{
		"payload": map[string]interface{}{
			"@a":      "puttype",
			"type":    "person",
			"indexes": []string{"name", "instrument"},
		},
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(&testDataType)
	if err != nil {
		t.Error(err)
		return
	}

	ty, err := tquery(t, cl, &buf)
	if err != nil {
		t.Error(err)
		return
	}

	if !ty.Success {
		t.Error(ty.Error)
		return
	}

	buf.Reset()
	err = json.NewEncoder(&buf).Encode(&testDataSet)
	if err != nil {
		t.Error(err)
		return
	}

	o, err := tquery(t, cl, &buf)
	if err != nil {
		t.Error(err)
		return
	}

	if !o.Success {
		t.Error(o.Error)
		return
	}

	reina := testDataSet["payload"]["data"].(map[string]interface{})
	m2, ok := o.Payload.(map[string]interface{})
	if !ok {
		t.Error("payload wasn't the right type")
	}

	err = assertMapEq(reina, m2)
	if err != nil {
		t.Error(err)
		return
	}

	testDataGet["payload"]["id"] = m2["id"]
	buf.Reset()
	err = json.NewEncoder(&buf).Encode(&testDataGet)
	if err != nil {
		t.Error(err)
		return
	}

	o, err = tquery(t, cl, &buf)
	if err != nil {
		t.Error(err)
		return
	}

	if !o.Success {
		t.Error(o.Error)
		return
	}

	m3, ok := o.Payload.([]interface{})[0].(map[string]interface{})
	if !ok {
		t.Error("payload wasn't the right type")
	}

	err = assertMapEq(reina, m3)
	if err != nil {
		t.Error(err)
		return
	}

	iml.Close()
}

func tquery(t *testing.T, cl *http.Client, buf *bytes.Buffer) (*OutgoingMessage, error) {

	rq, err := http.NewRequest("QUERY", "http://localhost/", buf)
	if err != nil {
		t.Error("request create error: ", err)
		return nil, err
	}

	rsp, err := cl.Do(rq)
	if err != nil {
		t.Error("request create error: ", err)
		return nil, err
	}

	var resp OutgoingMessage
	err = json.NewDecoder(rsp.Body).Decode(&resp)
	if err != nil {
		t.Error("decoder error: ", err)
		return nil, err
	}

	return &resp, nil
}

func TestBadQuery(t *testing.T) {
	iml, cl := getInmemoryServer()
	buf := bytes.NewBufferString("")
	msg, err := tquery(t, cl, buf)
	if err != nil {
		return
	}

	if msg.Success == true {
		t.Error("unsuccesfully successful")
		return
	}

	buf = bytes.NewBufferString(`{"payload":false}`)
	msg, err = tquery(t, cl, buf)
	if err != nil {
		return
	}

	if msg.Success == true {
		t.Error("unsuccesfully successful")
		return
	}

	iml.Close()
}

func assertMapEq(m1, m2 map[string]interface{}) error {
	for k, v := range m1 {
		if m2[k] != v {
			return fmt.Errorf("m2[%s] !== m1[%s]", k, k)
		}
	}

	return nil
}
