package debugger

import (
	"fmt"
	"io"
	"sync"
	"time"
)

type Event struct {
	At      time.Time
	Type    string
	Subject string
	Detail  string
}

type Tracer struct {
	mu     sync.Mutex
	events []Event
	out    io.Writer
}

func New(out io.Writer) *Tracer {
	return &Tracer{out: out, events: make([]Event, 0, 256)}
}

func (t *Tracer) Record(typ, subject, detail string) {
	e := Event{At: time.Now().UTC(), Type: typ, Subject: subject, Detail: detail}
	t.mu.Lock()
	t.events = append(t.events, e)
	t.mu.Unlock()
	if t.out != nil {
		fmt.Fprintf(t.out, "[%s] %-12s %-20s %s\n", e.At.Format(time.RFC3339Nano), e.Type, e.Subject, e.Detail)
	}
}

func (t *Tracer) Snapshot() []Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]Event, len(t.events))
	copy(out, t.events)
	return out
}
