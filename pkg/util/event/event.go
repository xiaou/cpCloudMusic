package event

import (
	"sync"
	"time"
)

type Event struct {
	c     chan struct{}
	isSet bool
	mutex sync.RWMutex
}

func NewEvent() *Event {
	res := &Event{}
	res.c = make(chan struct{}, 99)
	return res
}

func (e *Event) IsSet() bool {
	return e.isSet
}

func (e *Event) Set() {
	e.mutex.Lock()
	select {
	case e.c <- struct{}{}:
	default:
	}
	e.isSet = true
	e.mutex.Unlock()
}

func (e *Event) Clear() {
	e.mutex.Lock()
	for {
		shouldBreak := false
		select {
		case <-e.c:
			continue
		default:
			shouldBreak = true
		}
		if shouldBreak {
			break
		}
	}
	e.isSet = false
	e.mutex.Unlock()
}

func (e *Event) Wait() (signed bool) {
	<-e.c
	return true
}

func (e *Event) WaitUntil(timeout time.Duration) (signed bool) {
	select {
	case <-e.c:
		return true
	case <-time.After(timeout):
		return false
	}
}
