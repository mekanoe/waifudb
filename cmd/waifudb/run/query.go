package run

import (
	"github.com/kayteh/waifudb/db"
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
	case "puttype":
		s.queryPutType(ctx, p)
		return
	default:
		encodeOut(ctx, pktErrBadRequest)
	}

}

func (s *Server) queryGet(ctx *fasthttp.RequestCtx, p map[string]interface{}) {
	logger.WithField("p", p).Info("queryGet")
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

	// Mode 2: by index
	index, ok := p["index"].(string)
	if ok {
		val, ok := p["value"].(string)
		if !ok {
			encodeOut(ctx, pktErrBadRequest)
			return
		}

		i, err := s.w.GetItemByKey(ty, index, val)
		if err != nil {
			encodeOut(ctx, pktErrGetFailed)
			return
		}

		encodeOut(ctx, OutgoingMessage{
			Success: true,
			Payload: []map[string]interface{}{i},
		})

		return
	}

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

func (s *Server) queryPutType(ctx *fasthttp.RequestCtx, p map[string]interface{}) {
	ty := &db.Type{}

	var ok bool
	ty.Name, ok = p["type"].(string)
	if !ok {
		encodeOut(ctx, OutgoingMessage{
			Success: false,
			Error:   "`type` field required",
		})
		return
	}

	ifIdxs := p["indexes"].([]interface{})
	ty.Indexes = make([]string, len(ifIdxs))

	for k, v := range ifIdxs {
		ty.Indexes[k] = v.(string)
	}

	ifRels, _ := p["relations"].(map[string]interface{})
	ty.Relations = map[string]string{}

	for k, v := range ifRels {
		ty.Relations[k] = v.(string)
	}

	ty, err := s.w.CreateType(ty)
	if err != nil {
		encodeOut(ctx, OutgoingMessage{
			Success: false,
			Error:   "failed to create type",
		})
		return
	}

	encodeOut(ctx, OutgoingMessage{
		Success: true,
		Payload: ty,
	})
}
