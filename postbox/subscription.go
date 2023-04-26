package postbox

import (
	"fmt"
	"strings"

	"github.com/j13g/goutil/log"
	"github.com/j13g/goutil/types"
	"github.com/nats-io/nats.go"
)

func multiSubscribe(conn *nats.Conn, spec *SubjectSpec) (func(), chan *nats.Msg, error) {
	subjects := spec.Get()
	subscriptions := make([]*nats.Subscription, len(subjects))
	msgs := make(chan *nats.Msg)

	l := log.Get()
	for i, sub := range subjects {
		var err error
		l.Trace().Str("subject", sub).Msg("subscribing to subject")
		subscriptions[i], err = conn.ChanSubscribe(sub, msgs)
		if err != nil {
			panic(err) // TODO
		}
	}

	drainCallback := func() {
		for _, s := range subscriptions {
			s.Drain()
		}
	}

	return drainCallback, msgs, nil
}

func NewSubjectSpec() *SubjectSpec {
	return &SubjectSpec{}
}

type SubjectSpec struct {
	x []spec
}

func (s *SubjectSpec) AddPath(path string) *SubjectSpec {
	s.x = append(s.x, singleSpec(path))
	return s
}

func (s *SubjectSpec) AddPathF(path string, args ...any) *SubjectSpec {
	return s.AddPath(fmt.Sprintf(path, args...))
}

func (s *SubjectSpec) AddMulti(path string) *SubjectSpec {
	s.x = append(s.x, multiSpec{path: path})
	return s
}

func (s *SubjectSpec) AddMultiF(path string, args ...any) *SubjectSpec {
	return s.AddMulti(fmt.Sprintf(path, args...))
}

func (s *SubjectSpec) Exclude(paths ...string) *SubjectSpec {
	s.x = append(s.x, excludeSpec(paths))
	return s
}

func (s *SubjectSpec) Get() []string {
	paths := types.NewSet[string]()
	excludes := types.NewSet[string]()

	for _, spec := range s.x {
		switch spec.(type) {

		case excludeSpec:
			for _, exclude := range spec.get() {
				excludes.Add(exclude)
			}

		default:
			for _, path := range spec.get() {
				paths.Add(path)
			}
		}
	}

	result := paths.Subtract(excludes)
	return result.ToSlice()
}

type spec interface {
	get() []string
}

type excludeSpec []string

func (es excludeSpec) get() []string {
	return es
}

var _ spec = excludeSpec{}

type singleSpec string

func (s singleSpec) get() []string {
	return []string{string(s)}
}

type multiSpec struct {
	path string
}

func (s multiSpec) get() []string {
	path := ""
	set := types.NewSet[string]()
	parts := strings.Split(s.path, ".")
	for _, part := range parts {
		path += part
		set.Add(path)
		path += "."
	}
	return set.ToSlice()
}
