package postbox

import (
	"context"

	"github.com/j13g/goutil/log"
	"github.com/j13g/goutil/version"
	"github.com/rs/zerolog"
)

type Handler func(context.Context, *Envelope) *Envelope

func NewRouter() *Router {
	return &Router{
		x: make(map[string]*version.VersionFilterMap[Handler]),
		l: log.Get(),
	}
}

type Router struct {
	l zerolog.Logger
	x map[string]*version.VersionFilterMap[Handler]
}

func (r *Router) Register(name string, versionFilter string, h Handler) {
	if _, ok := r.x[name]; !ok {
		r.x[name] = version.NewVersionFilterMap[Handler]()
	}

	r.x[name].Add(versionFilter, h)
}

func (r *Router) Handle(ctx context.Context, req *Envelope) *Envelope {
	versionFilterMap, ok := r.x[req.Name]
	if !ok {
		r.l.Error().Str("name", req.Name).Msg("no handler found for message")
		return nil
	}

	handler := versionFilterMap.Get(req.Version)
	if handler.IsAbsent() {
		r.l.Error().
			Str("name", req.Name).
			Str("version", req.Version).
			Msg("no handler found for message version")
		return nil
	}

	response := handler.MustGet()(ctx, req)
	return response
}
