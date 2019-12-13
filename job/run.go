package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"time"

	jobEvent "gopkg.in/webnice/job.v1/event"
	jobTypes "gopkg.in/webnice/job.v1/types"
)

// Запуск процесса forkworker
func (jbo *impl) runForkWorker(prc *jobTypes.ForkWorker) (err error) {

	// TODO Скопировать запуск из wdaccessatomic

	//debug.Dumper(prc)

	return
}

// Запуск процесса worker
func (jbo *impl) runWorker(prc *jobTypes.Worker) (err error) {
	if prc.State.IsRun.Load().(bool) {
		err = jbo.Errors().ProcessAlreadyRunning()
		return
	}
	// Prepare()
	if err = safeCall(prc.Self.Prepare); err != nil {
		jbo.Event <- &jobEvent.Event{Act: jobEvent.EOnError, SourceID: prc.ID, Err: err}
		return
	}
	prc.State.IsRun.Store(true)
	jbo.Wg.Add(1)
	go func(prc *jobTypes.Worker) {
		defer func() {
			prc.State.IsRun.Store(false)
			safeWgDone(jbo.Wg)
			jbo.Event <- &jobEvent.Event{Act: jobEvent.EProcessStoped, SourceID: prc.ID}
		}()
		jbo.runProc(prc.Ctx, prc.Self, &prc.Pith)
	}(prc)

	return
}

// Запуск процесса task
func (jbo *impl) runTask(prc *jobTypes.Task) (err error) {
	if prc.State.IsRun.Load().(bool) {
		err = jbo.Errors().ProcessAlreadyRunning()
		return
	}
	prc.State.IsRun.Store(true)
	jbo.Wg.Add(1)
	go func(prc *jobTypes.Task) {
		defer func() {
			prc.State.IsRun.Store(false)
			safeWgDone(jbo.Wg)
			jbo.Event <- &jobEvent.Event{Act: jobEvent.EProcessStoped, SourceID: prc.ID}
		}()
		jbo.runProc(prc.Ctx, prc.Self, &prc.Pith)
	}(prc)

	return
}

// Запуск основного процесса с прерыванием
func (jbo *impl) runProc(ctx context.Context, pri jobTypes.BaseInterface, pith *jobTypes.Pith) {
	var pex = make(chan struct{})

	defer func() { close(pex) }()
	go func() {
		var err error
		defer func() { safeChannelSend(pex) }()
		jbo.Event <- &jobEvent.Event{Act: jobEvent.EProcessStarted, SourceID: pith.ID}
		err = safeCall(pri.Worker)
		if err != nil {
			// Отправка события ошибки
			jbo.Event <- &jobEvent.Event{Act: jobEvent.EOnError, SourceID: pith.ID, Err: err}
		}
		if err != nil && pith.State.Conf.Fatality {
			// Отправка события фатального завершения всех процессов и приложения
			jbo.Event <- &jobEvent.Event{Act: jobEvent.EProcessFatality, SourceID: pith.ID, Err: err}
		}
		if err == nil && pith.State.Conf.Restart {
			// Перезапуск остановившегося без ошибки процесса
			jbo.Event <- &jobEvent.Event{Act: jobEvent.ERestartProcess, SourceID: pith.ID, TargetID: pith.ID}
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
			jbo.Event <- &jobEvent.Event{Act: jobEvent.ECancelError, SourceID: pith.ID, Err: err}
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
