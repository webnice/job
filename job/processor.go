package job // import "gopkg.in/webnice/job.v1/job"

import "gopkg.in/webnice/debug.v1"
import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"time"

	"gopkg.in/webnice/job.v1/event"
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
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case prc.Task != nil:
			if prc.Task.State.IsRun.Load().(bool) {
				prc.Task.Cancel()
			}
		case prc.Worker != nil:
			if prc.Worker.State.IsRun.Load().(bool) {
				prc.Worker.Cancel()
			}
		case prc.ForkWorker != nil:
			if prc.ForkWorker.State.IsRun.Load().(bool) {
				prc.ForkWorker.Cancel()
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
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case prc.Task != nil:
			if prc.Task.ID != evt.TargetID {
				continue
			}
			// Таймаут перезапуска процесса
			if prc.Task.State.Conf.RestartTimeout > 0 {
				go func() {
					tmr := time.NewTimer(prc.Task.State.Conf.RestartTimeout)
					defer tmr.Stop()
					<-tmr.C
					jbo.Err = jbo.doTask(prc.Task.ID)
				}()
			} else {
				jbo.Err = jbo.doTask(prc.Task.ID)
			}
		case prc.Worker != nil:
			if prc.Worker.ID != evt.TargetID {
				continue
			}
			// Таймаут перезапуска процесса
			if prc.Worker.State.Conf.RestartTimeout > 0 {
				go func() {
					tmr := time.NewTimer(prc.Worker.State.Conf.RestartTimeout)
					defer tmr.Stop()
					<-tmr.C
					jbo.Err = jbo.doTask(prc.Worker.ID)
				}()
			} else {
				jbo.Err = jbo.doTask(prc.Worker.ID)
			}
		case prc.ForkWorker != nil:
			if prc.ForkWorker.ID != evt.TargetID {
				continue
			}
			if prc.ForkWorker.State.Conf.RestartTimeout > 0 {
				go func() {
					tmr := time.NewTimer(prc.ForkWorker.State.Conf.RestartTimeout)
					defer tmr.Stop()
					<-tmr.C
					jbo.Err = jbo.doTask(prc.ForkWorker.ID)
				}()
			} else {
				jbo.Err = jbo.doTask(prc.ForkWorker.ID)
			}
		}
	}
}

// Фатальная ошибка, остановка и выход
func (jbo *impl) eventFatality(evt *event.Event) {
	jbo.eventCancel()
}
