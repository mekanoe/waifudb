package run

import (
	"encoding/json"

	"net"

	"github.com/Sirupsen/logrus"
	"github.com/imdario/mergo"
	"github.com/kayteh/waifudb/datastore"
	"github.com/kayteh/waifudb/db"
	"github.com/valyala/fasthttp"
)

var (
	logger = logrus.WithFields(logrus.Fields{})

	pktErrBadRequest = OutgoingMessage{Success: false, Error: "bad request"}
	pktErrGetFailed  = OutgoingMessage{Success: false, Error: "couldn't get item"}
	pktErrNotFound   = OutgoingMessage{Success: false, Error: "item not found"}
	pktPong          = OutgoingMessage{Success: true, Payload: "pong"}
)

type OutgoingMessage struct {
	Success bool
	Payload interface{}
	Error   string
}

type IncomingMessage struct {
	// Type    string
	Payload interface{}
}

type Server struct {
	w   *db.WaifuDB
	cfg *Config
}

type Config struct {
	Addr   string
	DBPath string
}

func (c *Config) merge(incoming *Config) error {
	if incoming == nil {
		return nil
	}

	return mergo.MergeWithOverwrite(c, incoming)
}

var (
	defaultConfig = &Config{
		Addr:   "localhost:7099",
		DBPath: ".trash.db",
	}
)

func New(cfg *Config) (*Server, error) {
	c := defaultConfig
	err := c.merge(cfg)
	if err != nil {
		logger.WithError(err).Error("config merge error")
		return nil, err
	}

	st, err := datastore.New(&datastore.Config{
		Path: c.DBPath,
	})
	if err != nil {
		logger.WithError(err).Error("failed to get datastore")
		return nil, err
	}

	w, err := db.New(st)
	if err != nil {
		logger.WithError(err).Error("failed to rev up db")
		return nil, err
	}

	s := &Server{
		w:   w,
		cfg: c,
	}

	return s, nil
}

func (s *Server) getServer() *fasthttp.Server {
	return &fasthttp.Server{
		Name:    "waifudb",
		Handler: s.handler,
	}
}

func (s *Server) Serve(n net.Listener) {
	srv := s.getServer()
	srv.Serve(n)
}

func (s *Server) ListenAndServe() {
	srv := s.getServer()
	srv.ListenAndServe(s.cfg.Addr)
}

func encodeOut(ctx *fasthttp.RequestCtx, data OutgoingMessage) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.WithError(err).Error("encode out: marshaling failure")
	}

	ctx.Write(body)
}

func decodeIn(ctx *fasthttp.RequestCtx) (IncomingMessage, error) {
	var i IncomingMessage
	err := json.Unmarshal(ctx.PostBody(), &i)
	return i, err
}

func (s *Server) handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "QUERY":
		s.queryHandler(ctx)
		return
	case "PING":
		s.pingHandler(ctx)
		return
	default:
		encodeOut(ctx, pktErrBadRequest)
		return
	}
}

func (s *Server) pingHandler(ctx *fasthttp.RequestCtx) {
	encodeOut(ctx, pktPong)
}
