package infra

import (
	"reflect"
	"time"
)

type event = interface{}
type message struct {
	stream    string
	name      string
	value     event
	createdAt time.Time
}

// Events - actual storage for all events from all Aggregate's
//
// Place where ACID is introduced.
// No resilient implementation.
// I would use Optimistic Concurrency Control algorithm in order to make every write ACID.
// Additional in-memory cache is required in order to avoid loading all events from storage on every read.
type Events map[string][]message

func (r Events) read(a Aggregate) error {
	for _, m := range r[a.ID()] {
		if err := a.Commit(m.value, m.createdAt); err != nil {
			return err
		}
	}

	return nil
}

func (r Events) write(a Aggregate) error {
	n := time.Now()
	s := a.ID()
	for _, e := range a.Uncommitted(true) {
		m := message{
			stream:    s,
			value:     e,
			name:      reflect.TypeOf(e).Name(),
			createdAt: time.Now(),
		}
		r[s] = append(r[s], m)

		if err := a.Commit(e, n); err != nil {
			return err
		}

		log("DBG #%s|%s", m.stream, m.name)
	}

	return nil
}

type Aggregate interface {
	ID() string
	Uncommitted(bool) []event
	Commit(event, time.Time) error
}

var log = DefaultLogger.Tag("EventStore").Print
