package event

import (
	"log"
	"testing"
	"time"
)

var (
	e = NewEvent()
)

func TestSend(t *testing.T) {
	go func() {
		for {
			e.Set()
			e.Set()
			time.Sleep(time.Second * 3)
		}
	}()
}

func TestWait(t *testing.T) {
	res := e.Wait()
	log.Printf("Wait return %v, then Cleared.\n", res)
	e.Clear()
}

func TestWaitUntil(t *testing.T) {
	go func() {
		for {
			res := e.WaitUntil(time.Second * 2)
			log.Printf("WaitUntil return %v\n", res)
		}
	}()
}

func TestXXX(t *testing.T) {
	time.Sleep(time.Second * 99)
}
