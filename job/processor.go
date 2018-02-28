package job // import "gopkg.in/webnice/job.v1/job"

import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"time"

	"gopkg.in/webnice/job.v1/event"
	"gopkg.in/webnice/job.v1/types"
)

// Горутина обработки событий
func (jbo *impl) EventProcessor() {
	var evt *event.Event

	for {
		// nil приходит при пересоздании канала по команде Reset()
		if evt = <-jbo.Event; evt == nil {
			continue
		}
		switch evt.Act {
		// Отправка всем запущенным процессам сигнала Cancel()
		case event.ECancel:
			jbo.eventCancel()
		// События ошибок
		case event.EOnError, event.ECancelError:
			jbo.eventError(evt)
		// Изменение состояния процесса
		case event.EProcessStarted:
			jbo.eventChangeState(evt, true)
		case event.EProcessStoped:
			jbo.eventChangeState(evt, false)
		case event.ERestartProcess:
			jbo.eventRestartProcess(evt)
		case event.EProcessFatality:
			jbo.eventFatality(evt)

		// Любое не известное событие
		default:
			log.Debugf("\n%s", debug.DumperString(evt))
		}
	}
}

// Сигнал завершения всех запущенных процессов
func (jbo *impl) eventCancel() {
	var elm *list.Element
	var prc *Process
	var ok bool

	jbo.Exit.Store(true)
	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok || prc == nil {
			continue
		}
		switch wrk := prc.P.(type) {
		case *types.Task:
			if wrk.State.IsRun.Load().(bool) {
				wrk.Cancel()
			}
		case *types.Worker:
			if wrk.State.IsRun.Load().(bool) {
				wrk.Cancel()
			}
		case *types.ForkWorker:
			if wrk.State.IsRun.Load().(bool) {
				wrk.Cancel()
			}
		}
	}
}

// Выполнение внешней функции принимающей события ошибок
func (jbo *impl) eventError(evt *event.Event) {
	defer func() { _ = recover() }()

	if jbo.ErrorFunc == nil {
		return
	}
	jbo.ErrorFunc(evt.SourceID, evt.Err)
}

// Выполнение внешней функции принимающей событие изменения статуса процесса
func (jbo *impl) eventChangeState(evt *event.Event, running bool) {
	defer func() { _ = recover() }()

	if jbo.ChangeStateFunc == nil {
		return
	}
	jbo.ChangeStateFunc(evt.SourceID, running)
}

// Событие перезапуска процесса завершившегося без ошибки
func (jbo *impl) eventRestartProcess(evt *event.Event) {
	var elm *list.Element
	var prc *Process
	var ok bool

	if jbo.Exit.Load().(bool) {
		return
	}
	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok || prc == nil {
			continue
		}
		switch wrk := prc.P.(type) {
		case *types.Task:
			if wrk.ID != evt.TargetID {
				continue
			}
			// Перезапуск процесса с таймаутом
			go jbo.doTaskWithTimeout(wrk, wrk.State.Conf.RestartTimeout)
		case *types.Worker:
			if wrk.ID != evt.TargetID {
				continue
			}
			// Перезапуск процесса с таймаутом
			go jbo.doTaskWithTimeout(wrk, wrk.State.Conf.RestartTimeout)
		case *types.ForkWorker:
			if wrk.ID != evt.TargetID {
				continue
			}
			// Перезапуск процесса с таймаутом
			go jbo.doTaskWithTimeout(wrk, wrk.State.Conf.RestartTimeout)
		}
	}
}

// Перезапуск процесса с таймаутом
func (jbo *impl) doTaskWithTimeout(proc interface{}, tm time.Duration) {
	var err error
	var tmr *time.Timer
	var id string

	if tm > 0 {
		tmr = time.NewTimer(tm)
		defer tmr.Stop()
		<-tmr.C
	}
	switch proc.(type) {
	case *types.Task:
		if item, ok := proc.(*types.Task); ok {
			id, err = item.ID, jbo.runTask(item)
		}
	case *types.Worker:
		if item, ok := proc.(*types.Worker); ok {
			id, err = item.ID, jbo.runWorker(item)
		}
	case *types.ForkWorker:
		if item, ok := proc.(*types.ForkWorker); ok {
			id, err = item.ID, jbo.runForkWorker(item)
		}
	}
	if err != nil {
		jbo.Event <- &event.Event{Act: event.EOnError, SourceID: id, Err: err}
	}
}

// Фатальная ошибка, остановка и выход
func (jbo *impl) eventFatality(evt *event.Event) {
	jbo.eventCancel()
}
