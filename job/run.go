package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"time"

	"gopkg.in/webnice/job.v1/event"
	"gopkg.in/webnice/job.v1/types"
)

// Запуск процесса forkworker
func (jbo *impl) runForkWorker(prc *types.ForkWorker) (err error) {

	// TODO
	// Скопировать запуск из wdaccessatomic

	//debug.Dumper(prc)

	return
}

// Запуск процесса worker
func (jbo *impl) runWorker(prc *types.Worker) (err error) {
	if prc.State.IsRun.Load().(bool) {
		err = ErrorProcessAlreadyRunning()
		return
	}
	// Prepare()
	if err = safeCall(prc.Self.Prepare); err != nil {
		jbo.Event <- &event.Event{Act: event.EOnError, SourceID: prc.ID, Err: err}
		return
	}
	prc.State.IsRun.Store(true)
	jbo.Wg.Add(1)
	go func(prc *types.Worker) {
		defer func() {
			prc.State.IsRun.Store(false)
			safeWgDone(jbo.Wg)
			jbo.Event <- &event.Event{Act: event.EProcessStoped, SourceID: prc.ID}
		}()
		jbo.runProc(prc.Ctx, prc.Self, &prc.Pith)
	}(prc)

	return
}

// Запуск процесса task
func (jbo *impl) runTask(prc *types.Task) (err error) {
	if prc.State.IsRun.Load().(bool) {
		err = ErrorProcessAlreadyRunning()
		return
	}
	prc.State.IsRun.Store(true)
	jbo.Wg.Add(1)
	go func(prc *types.Task) {
		defer func() {
			prc.State.IsRun.Store(false)
			safeWgDone(jbo.Wg)
			jbo.Event <- &event.Event{Act: event.EProcessStoped, SourceID: prc.ID}
		}()
		jbo.runProc(prc.Ctx, prc.Self, &prc.Pith)
	}(prc)

	return
}

// Запуск основного процесса с прерыванием
func (jbo *impl) runProc(ctx context.Context, pri types.BaseInterface, pith *types.Pith) {
	var pex = make(chan struct{})

	defer func() { close(pex) }()
	go func() {
		var err error
		defer func() { safeChannelSend(pex) }()
		jbo.Event <- &event.Event{Act: event.EProcessStarted, SourceID: pith.ID}
		err = safeCall(pri.Worker)
		if err != nil {
			// Отправка события ошибки
			jbo.Event <- &event.Event{Act: event.EOnError, SourceID: pith.ID, Err: err}
		}
		if err != nil && pith.State.Conf.Fatality {
			// Отправка события фатального завершения всех процессов и приложения
			jbo.Event <- &event.Event{Act: event.EProcessFatality, SourceID: pith.ID, Err: err}
		}
		if err == nil && pith.State.Conf.Restart {
			// Перезапуск остановившегося без ошибки процесса
			jbo.Event <- &event.Event{Act: event.ERestartProcess, SourceID: pith.ID, TargetID: pith.ID}
		}
	}()
	select {
	// Завершение выполнения процесса
	case <-pex:
		return
	// Выполнение функции Cancel()
	case <-ctx.Done():
		// Прерывание выполнения
		if err := safeCall(pri.Cancel); err != nil {
			jbo.Event <- &event.Event{Act: event.ECancelError, SourceID: pith.ID, Err: err}
		}
		if pith.State.Conf.CancelTimeout > 0 {
			// Таймаут ожидания завершения процесса после выполнения функции Cancel()
			// После этого ожидания, отправляем в канал сигнал, как буд-то процесс завершился
			go func() {
				if pith.State.Conf.CancelTimeout > 0 {
					tmr := time.NewTimer(pith.State.Conf.CancelTimeout)
					defer tmr.Stop()
					<-tmr.C
				}
				safeChannelSend(pex)
			}()
		}
	}
	// Ожидание завершения процесса
	<-pex
}
