package run

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

var (
	logger = logrus.WithFields(logrus.Fields{})

	pktErrBadRequest = OutgoingMessage{Success: false, Payload: "bad request"}
	pktPong          = OutgoingMessage{Success: true, Payload: "pong"}
)

type OutgoingMessage struct {
	Success bool
	Payload interface{}
}

type IncomingMessage struct {
	Type    string
	Payload interface{}
}

func Start(addr string) {
	server := &fasthttp.Server{
		Name:    "waifudb",
		Handler: handler,
	}

	// TODO: make configurable
	err := server.ListenAndServe(addr)
	if err != nil {
		logger.WithError(err).Fatal("listen failed")
	}
}

func encodeOut(ctx *fasthttp.RequestCtx, data OutgoingMessage) {
	body, err := json.Marshal(data)
	if err != nil {
		logger.WithError(err).Error("encode out: marshaling failure")
	}

	ctx.Write(body)
}

func handler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Method()) {
	case "QUERY":
		queryHandler(ctx)
		return
	case "PING":
		pingHandler(ctx)
		return
	default:
		encodeOut(ctx, pktErrBadRequest)
		return
	}
}

func pingHandler(ctx *fasthttp.RequestCtx) {
	encodeOut(ctx, pktPong)
}
