package job

import jobTypes "github.com/webnice/job/v2/types"

// IsRun Для процессов в состоянии запущен возвращается истина.
func (prc *Process) IsRun() (ret bool, err error) {
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		ret = wrk.State.IsRun.Load().(bool)
	case *jobTypes.Worker:
		ret = wrk.State.IsRun.Load().(bool)
	case *jobTypes.ForkWorker:
		ret = wrk.State.IsRun.Load().(bool)
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}

// ID Возвращается идентификатор процесса.
func (prc *Process) ID() (ret string, err error) {
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		ret = wrk.ID
	case *jobTypes.Worker:
		ret = wrk.ID
	case *jobTypes.ForkWorker:
		ret = wrk.ID
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}

// State Возвращается указатель на состояние процесса.
func (prc *Process) State() (ret *jobTypes.State, err error) {
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		ret = &wrk.State
	case *jobTypes.Worker:
		ret = &wrk.State
	case *jobTypes.ForkWorker:
		ret = &wrk.State
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}

// InfoRequest Запрос конфигурации процесса.
func (prc *Process) InfoRequest() (ret *jobTypes.Configuration, err error) {
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		ret = wrk.Self.Info(wrk.ID)
	case *jobTypes.Worker:
		ret = wrk.Self.Info(wrk.ID)
	case *jobTypes.ForkWorker:
		ret = wrk.Self.Info(wrk.ID)
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}

// Configuration Установка конфигурации процессу.
func (prc *Process) Configuration(cfg *jobTypes.Configuration) (err error) {
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		wrk.State.Conf = cfg
	case *jobTypes.Worker:
		wrk.State.Conf = cfg
	case *jobTypes.ForkWorker:
		wrk.State.Conf = cfg
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}
