package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"time"
)

// Constructor of singleton
func init() {
	singleton = new(impl)
	singleton.Reset()
	go singleton.EventProcessor()
}

// Error Последняя внутренняя ошибка
func (jbo *impl) Error() error { return jbo.Err }

// RegisterErrorFunc Регистрация функции получения ошибок в работе управляемых процессов
func (jbo *impl) RegisterErrorFunc(fn OnErrorFunc) Interface { jbo.ErrorFunc = fn; return jbo }

// RegisterChangeStateFunc Регистрация функции получения изменения состояния процессов
func (jbo *impl) RegisterChangeStateFunc(fn OnChangeStateFunc) Interface {
	jbo.ChangeStateFunc = fn
	return jbo
}

// UnregisterErrorFunc Удаление ранее зарегистрированной функции получения ошибок
func (jbo *impl) UnregisterErrorFunc() Interface { jbo.ErrorFunc = nil; return jbo }

// Wait Ожидание завершения всех работающих процессов
func (jbo *impl) Wait() Interface {
	var wgc = make(chan struct{})
	defer func() { close(wgc); safeWgDoneForAll(jbo.Wg) }()
	go func() { defer func() { safeChannelSend(wgc) }(); safeWgWait(jbo.Wg) }()
	select {
	case <-wgc:
		jbo.Err = nil
	case <-jbo.Ctx.Done():
		switch jbo.Ctx.Err() {
		case context.Canceled:
			jbo.Err = ErrorProcessesAreStillRunning()
		case context.DeadlineExceeded:
			jbo.Err = ErrorDeadlineExceeded()
		}
	}
	return jbo
}

// WaitWithTimeout Ожидание завершения всех работающих процессов, но не более чем время указанное в timeout
func (jbo *impl) WaitWithTimeout(timeout time.Duration) Interface {
	var wgc = make(chan struct{})
	jbo.Ctx, jbo.CancelFunc = context.WithTimeout(jbo.Ctx, timeout)
	defer func() { close(wgc) }()
	go func() { defer func() { safeChannelSend(wgc) }(); jbo.Wait() }()
	select {
	case <-wgc:
		jbo.CancelFunc()
	case <-jbo.Ctx.Done():
		if jbo.Ctx.Err() == context.DeadlineExceeded {
			jbo.Err = ErrorDeadlineExceeded()
		}
	}
	return jbo
}
