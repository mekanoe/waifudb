package run

import (
	"github.com/valyala/fasthttp"
)

func (s *Server) queryHandler(ctx *fasthttp.RequestCtx) {
	m, err := decodeIn(ctx)
	if err != nil {
		encodeOut(ctx, OutgoingMessage{Success: false, Error: err.Error()})
		return
	}

	p, ok := m.Payload.(map[string]interface{})
	if !ok {
		encodeOut(ctx, pktErrBadRequest)
		logger.WithField("p", m.Payload).Error("payload wasn't a valid query...")
		return
	}

	switch p["@a"] {
	case "get":
		s.queryGet(ctx, p)
		return
	case "set":
		s.querySet(ctx, p)
		return
	default:
		encodeOut(ctx, pktErrBadRequest)
	}

}

func (s *Server) queryGet(ctx *fasthttp.RequestCtx, p map[string]interface{}) {
	ty, ok := p["type"].(string)
	if !ok {
		encodeOut(ctx, pktErrBadRequest)
	}

	// Mode 1: by ID
	id, ok := p["id"].(string)
	if ok {
		i, err := s.w.GetItem(ty, id)
		if err != nil {
			encodeOut(ctx, pktErrGetFailed)
		}

		encodeOut(ctx, OutgoingMessage{
			Success: true,
			Payload: []map[string]interface{}{i},
		})
		return
	}
	/*
		// Mode 2: by index
		index, ok := p["index"].(string)
		if ok {
			val, ok := p["value"].(string)
			if !ok {
				encodeOut(ctx, pktErrBadRequest)
				return
			}
			return
		}*/

	encodeOut(ctx, pktErrBadRequest)
}

func (s *Server) querySet(ctx *fasthttp.RequestCtx, p map[string]interface{}) {
	ty, ok := p["type"].(string)
	if !ok {
		encodeOut(ctx, pktErrBadRequest)
		return
	}

	data, ok := p["data"].(map[string]interface{})
	if !ok {
		encodeOut(ctx, pktErrBadRequest)
		return
	}

	i, err := s.w.PutItem(ty, data)
	if err != nil {
		logger.WithError(err).Error("putitem failed")
		encodeOut(ctx, OutgoingMessage{
			Success: false,
			Error:   "failed to write",
		})
		return
	}

	encodeOut(ctx, OutgoingMessage{
		Success: true,
		Payload: i,
	})
}
