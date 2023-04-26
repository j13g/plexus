package postbox

import (
	"context"
	"github.com/j13g/goutil/log"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"sync"
	"time"
)

func NewInbox(conn *nats.Conn, appName, appVersion, nodeName string) *Inbox {
	return &Inbox{
		l:      log.Get(),
		conn:   conn,
		router: NewRouter(),

		app:        appName,
		appVersion: appVersion,
		nodeName:   nodeName,
	}
}

type Inbox struct {
	l      zerolog.Logger
	router *Router
	conn   *nats.Conn

	nodeName   string
	app        string
	appVersion string

	wg   sync.WaitGroup
	stop chan struct{}
}

func (i *Inbox) Router() *Router {
	return i.router
}

func (i *Inbox) Start(subjectSpec *SubjectSpec) {
	i.stop = make(chan struct{})
	i.l.Debug().Msg("starting message server worker")

	i.wg.Add(1)
	go i.worker(subjectSpec)
}

func (i *Inbox) Stop() {
	close(i.stop)
	i.wg.Wait()
	i.l.Debug().Msg("message server worker stopped")
}

func (i *Inbox) worker(subjectSpec *SubjectSpec) {
	defer i.wg.Done()

	drainFunc, msgsChan, err := multiSubscribe(i.conn, subjectSpec)
	if err != nil {
		return
	}
	defer drainFunc()

	for {
		select {
		case <-i.stop:
			i.l.Trace().Msg("got shutdown signal")
			return
		case natsMsg := <-msgsChan:
			start := time.Now().UTC()
			request, err := readEnvelope(natsMsg.Data)
			if err != nil {
				i.l.Error().Err(err).Msg("failed to parse request from message")
				continue
			}
			i.l.Trace().Interface("request", request).Msg("recieved message from NATS")

			response := i.router.Handle(context.Background(), request)

			if response.Meta == nil {
				response.Meta = make(map[string]any)
			}
			response.TS = time.Now().UTC()
			response.Meta["processing_time"] = response.TS.Sub(start).String()
			response.Meta["node_name"] = i.nodeName
			response.Meta["app"] = i.app
			response.Meta["app_version"] = i.appVersion

			data, err := writeEnvelope(response)
			if err != nil {
				panic(err) // TODO
			}

			if natsMsg.Reply != "" {
				i.l.Trace().Interface("response", response).Msg("sending response")
				err := natsMsg.Respond(data)
				if err != nil {
					panic(err) // TODO
				}
			}
		}
	}
}
