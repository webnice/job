package job

import (
	"fmt"
	runtimeDebug "runtime/debug"
	"sync"
)

// Безопасный метод отправки пустой структуры (сигнала) в канал.
func safeChannelSend(ch chan<- struct{}) {
	defer func() { _ = recover() }()
	ch <- struct{}{}
}

// Безопасно делает один вызов Done() для sync.WaitGroup.
func safeWgDone(wg *sync.WaitGroup) {
	defer func() { _ = recover() }()
	wg.Done()
}

// Безопасно обнуляет до конца sync.WaitGroup.
func safeWgDoneForAll(wg *sync.WaitGroup) {
	defer func() { _ = recover() }()
	for {
		wg.Done()
	}
}

// Безопасно выполняет Wait() для sync.WaitGroup.
func safeWgWait(wg *sync.WaitGroup) {
	defer func() { _ = recover() }()
	wg.Wait()
}

// Безопасный запуск функции.
func safeCall(fn func() error) (err error) {
	var ok bool

	defer func() {
		if e := recover(); e != nil {
			if err, ok = e.(error); !ok {
				err = fmt.Errorf("%v", e)
			}
			err = fmt.Errorf("%s\n%s", err, string(runtimeDebug.Stack()))
		}
	}()
	err = fn()

	return
}
