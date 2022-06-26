package job

import (
	"container/list"
	"log"
	"time"

	jobEvent "github.com/webnice/job/v2/event"
	jobTypes "github.com/webnice/job/v2/types"
)

// EventProcessor Горутина обработки событий.
func (jbo *impl) EventProcessor() {
	var evt *jobEvent.Event

	for {
		// nil приходит при пересоздании канала по команде Reset().
		if evt = <-jbo.Event; evt == nil {
			continue
		}
		switch evt.Act {
		// Отправка всем запущенным процессам сигнала Cancel().
		case jobEvent.ECancel:
			jbo.eventCancel()
		// События ошибок.
		case jobEvent.EOnError, jobEvent.ECancelError:
			jbo.eventError(evt)
		// Изменение состояния процесса.
		case jobEvent.EProcessStarted:
			jbo.eventChangeState(evt, true)
		case jobEvent.EProcessStopped:
			jbo.eventChangeState(evt, false)
		case jobEvent.ERestartProcess:
			jbo.eventRestartProcess(evt)
		case jobEvent.EProcessFatality:
			jbo.eventFatality(evt)
		// Любое не известное событие.
		default:
			log.Printf("not implemented event: %q", string(evt.Act))
		}
	}
}

// Сигнал завершения всех запущенных процессов.
func (jbo *impl) eventCancel() {
	var (
		elm *list.Element
		prc *Process
		ok  bool
	)

	jbo.Exit.Store(true)
	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok || prc == nil {
			continue
		}
		switch wrk := prc.P.(type) {
		case *jobTypes.Task:
			if wrk.State.IsRun.Load().(bool) {
				wrk.Cancel()
			}
		case *jobTypes.Worker:
			if wrk.State.IsRun.Load().(bool) {
				wrk.Cancel()
			}
		case *jobTypes.ForkWorker:
			if wrk.State.IsRun.Load().(bool) {
				wrk.Cancel()
			}
		}
	}
}

// Выполнение внешней функции принимающей события ошибок.
func (jbo *impl) eventError(evt *jobEvent.Event) {
	defer func() { _ = recover() }()

	if jbo.ErrorFunc == nil {
		return
	}
	jbo.ErrorFunc(evt.SourceID, evt.Err)
}

// Выполнение внешней функции принимающей событие изменения статуса процесса.
func (jbo *impl) eventChangeState(evt *jobEvent.Event, running bool) {
	defer func() { _ = recover() }()

	if jbo.ChangeStateFunc == nil {
		return
	}
	jbo.ChangeStateFunc(evt.SourceID, running)
}

// Событие перезапуска процесса завершившегося без ошибки.
func (jbo *impl) eventRestartProcess(evt *jobEvent.Event) {
	if jbo.Exit.Load().(bool) {
		return
	}

	jbo.err = jbo.RegisteredProcessIterate(func(elm *list.Element, prc *Process) (e error) {
		switch wrk := prc.P.(type) {
		case *jobTypes.Task:
			if wrk.ID != evt.TargetID {
				return
			}
			// Перезапуск процесса с таймаутом.
			go jbo.doTaskWithTimeout(wrk, wrk.State.Conf.RestartTimeout)
		case *jobTypes.Worker:
			if wrk.ID != evt.TargetID {
				return
			}
			// Перезапуск процесса с таймаутом.
			go jbo.doTaskWithTimeout(wrk, wrk.State.Conf.RestartTimeout)
		case *jobTypes.ForkWorker:
			if wrk.ID != evt.TargetID {
				return
			}
			// Перезапуск процесса с таймаутом.
			go jbo.doTaskWithTimeout(wrk, wrk.State.Conf.RestartTimeout)
		}

		return
	})
}

// Перезапуск процесса с таймаутом.
func (jbo *impl) doTaskWithTimeout(proc interface{}, tm time.Duration) {
	var (
		err error
		tmr *time.Timer
		id  string
	)

	if tm > 0 {
		tmr = time.NewTimer(tm)
		defer tmr.Stop()
		<-tmr.C
	}
	switch proc.(type) {
	case *jobTypes.Task:
		if item, ok := proc.(*jobTypes.Task); ok {
			id, err = item.ID, jbo.runTask(item)
		}
	case *jobTypes.Worker:
		if item, ok := proc.(*jobTypes.Worker); ok {
			id, err = item.ID, jbo.runWorker(item)
		}
	case *jobTypes.ForkWorker:
		if item, ok := proc.(*jobTypes.ForkWorker); ok {
			id, err = item.ID, jbo.runForkWorker(item)
		}
	}
	if err != nil {
		jbo.Event <- &jobEvent.Event{Act: jobEvent.EOnError, SourceID: id, Err: err}
	}
}

// Фатальная ошибка, остановка и выход.
func (jbo *impl) eventFatality(evt *jobEvent.Event) {
	jbo.eventCancel()
}
