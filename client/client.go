// Package client is a reference client for waifudb
package client

import (
	"encoding/json"
	"errors"
	"net/http"

	"bytes"

	"github.com/imdario/mergo"
)

var (
	ErrDatabaseNotReady = errors.New("waifudb/client: database not ready")
	ErrGeneralError     = errors.New("waifudb/client: general error 🤷‍")
)

// Client is just a regular WaifuDB client!
// It must be created with New() to prevent oddities.
type Client struct {
	cfg       *Config
	connector *http.Client
}

// Config contains any values needed to make client requests
//
// - Addr: HTTP addr to the database or cluster (default: `http://localhost:9077`)
type Config struct {
	Addr string
}

func (c *Config) merge(incoming *Config) error {
	if incoming == nil {
		return nil
	}

	return mergo.MergeWithOverwrite(c, incoming)
}

var defaultConfig = &Config{
	Addr: "http://localhost:7099/",
}

// New creates a WaifuDB Client, ready to use immediately!
func New(cfg *Config) (*Client, error) {
	c := defaultConfig

	err := c.merge(cfg)
	if err != nil {
		return nil, err
	}

	cl := &Client{
		cfg:       c,
		connector: &http.Client{},
	}

	ok, err := cl.Ping()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrDatabaseNotReady
	}

	return cl, nil
}

func (cl *Client) request(method string, data interface{}) (*http.Response, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		return nil, err
	}

	rq, err := http.NewRequest(method, cl.cfg.Addr, &buf)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(rq)
}

func fromJSON(rsp *http.Response, ptr interface{}) error {
	return json.NewDecoder(rsp.Body).Decode(ptr)
}

// Ping the database to check for aliveness
func (cl *Client) Ping() (bool, error) {
	rsp, err := cl.request("PING", map[string]interface{}{"ping": 1})
	if err != nil {
		return false, err
	}

	if rsp.StatusCode != 200 {
		return false, ErrDatabaseNotReady
	}

	var pong map[string]interface{}
	err = fromJSON(rsp, &pong)
	if err != nil {
		return false, err
	}

	if pong["Success"].(bool) && pong["Payload"].(string) == "pong" {
		return true, nil
	}

	return false, ErrGeneralError
}
