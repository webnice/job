package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"context"
	"time"

	jobTypes "gopkg.in/webnice/job.v1/types"
)

// Constructor of singleton
func init() {
	singleton = new(impl)
	singleton.Reset()
	go singleton.EventProcessor()
}

// Error Последняя внутренняя ошибка
func (jbo *impl) Err() error { return jbo.err }

// Errors Ошибки известного состояни, которые могут вернуть функции пакета
func (jbo *impl) Errors() *Error { return Errors() }

// RegisterErrorFunc Регистрация функции получения ошибок о работе управляемых процессов
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
		jbo.err = nil
	case <-jbo.Ctx.Done():
		switch jbo.Ctx.Err() {
		case context.Canceled:
			jbo.err = jbo.Errors().ProcessesAreStillRunning()
		case context.DeadlineExceeded:
			jbo.err = jbo.Errors().DeadlineExceeded()
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
			jbo.err = jbo.Errors().DeadlineExceeded()
		}
	}

	return jbo
}

// RegisteredProcessIterate Итерация по всем зарегистрированным процессам
func (jbo *impl) RegisteredProcessIterate(fn func(*list.Element, *Process) error) (err error) {
	var (
		elm *list.Element
		prc *Process
		ok  bool
	)

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok {
			err = jbo.Errors().TypeNotImplemented()
			return
		}
		if err = fn(elm, prc); err != nil {
			return
		}
	}

	return
}

// RegisteredProcessFindByID Поиск зарегистрированного процесса по ID
func (jbo *impl) RegisteredProcessFindByID(id string) (item *list.Element, ret *Process, err error) {
	var (
		found bool
	)

	if err = jbo.RegisteredProcessIterate(
		func(elm *list.Element, prc *Process) (e error) {
			var (
				src  string
				full bool
			)
			if src, e = prc.ID(); e != nil {
				return
			}
			if full, _ = jbo.compareID(src, id, 0); full {
				item, ret, found = elm, prc, true
			}
			return
		}); err != nil {
		return
	}
	if !found {
		err = jbo.Errors().ProcessNotFound()
		return
	}

	return
}

// ProcessObjectReturnToPool Возвращение объекта процесса в пул
func (jbo *impl) ProcessObjectReturnToPool(prc *Process) (err error) {
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		jbo.Pool.TaskPut(wrk)
	case *jobTypes.Worker:
		jbo.Pool.WorkerPut(wrk)
	case *jobTypes.ForkWorker:
		jbo.Pool.ForkWorkerPut(wrk)
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}
