package job

import (
	"runtime"
	"sync"
	"testing"
)

func TestSafeChannelSend(t *testing.T) {
	var tch chan struct{}

	defer func() {
		if e := recover(); e != nil {
			t.Errorf("Error in safeChannelSend()")
		}
	}()

	tch = make(chan struct{})
	go func() { <-tch }()
	safeChannelSend(tch)

	close(tch)
	safeChannelSend(tch)
}

func TestSafeWgDone(t *testing.T) {
	var twg *sync.WaitGroup
	var waitOk bool

	defer func() {
		if e := recover(); e != nil {
			t.Errorf("Error in safeWgDone()")
		}
	}()

	twg = new(sync.WaitGroup)
	twg.Add(1)
	go func(wg *sync.WaitGroup) { safeWgWait(twg); waitOk = true }(twg)

	safeWgDone(twg)
	runtime.Gosched()
	safeWgDone(twg)
	runtime.Gosched()
	safeWgDone(twg)
	if !waitOk {
		t.Errorf("Error in safeWgWait()")
	}
}

func TestSafeCall(t *testing.T) {
	var err error
	var f func() error

	defer func() {
		if e := recover(); e != nil {
			t.Errorf("Error in safeCall()")
		}
	}()

	f = func() error {
		panic("OK")
		return nil
	}
	if err = safeCall(f); err == nil {
		t.Errorf("Error in safeCall()")
	}
}
