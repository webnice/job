package job

import (
	"container/list"
	"sort"

	jobEvent "github.com/webnice/job/event"
	jobTypes "github.com/webnice/job/types"
)

// Запрос конфигурации всех процессов
func (jbo *impl) getConfiguration() (err error) {
	var (
		cfg *jobTypes.Configuration
	)

	if err = jbo.RegisteredProcessIterate(
		func(elm *list.Element, prc *Process) (e error) {
			if cfg, e = prc.InfoRequest(); e != nil {
				return
			}
			if cfg == nil {
				cfg = jobTypes.DefaultConfiguration()
			}
			e = prc.Configuration(cfg)
			return
		}); err != nil {
		return
	}

	return
}

// Составление списка процессов в соответствии с приоритетами запуска и остановки
func (jbo *impl) priority() {
	type plist struct {
		Start int32
		Stop  int32
		ID    string
	}
	var (
		pst []*plist
		n   int
	)

	if jbo.err = jbo.RegisteredProcessIterate(
		func(elm *list.Element, prc *Process) (e error) {
			var (
				id string
				st *jobTypes.State
			)
			if id, e = prc.ID(); e != nil {
				return
			}
			if st, e = prc.State(); e != nil {
				return
			}
			pst = append(pst, &plist{ID: id, Start: st.Conf.PriorityStart, Stop: st.Conf.PriorityStop})
			return
		}); jbo.err != nil {
		return
	}
	// В порядке старта
	sort.Slice(pst, func(i int, j int) bool { return pst[i].Start < pst[j].Start })
	jbo.StartPriority = make([]string, len(pst))
	for n = range pst {
		jbo.StartPriority[n] = pst[n].ID
	}
	// В порядке остановки
	sort.Slice(pst, func(i int, j int) bool { return pst[i].Stop < pst[j].Stop })
	jbo.StopPriority = make([]string, len(pst))
	for n = range pst {
		jbo.StopPriority[n] = pst[n].ID
	}
}

// Do Запуск библиотеки, подготовка и запуск процессов с флагом Autostart
// Ошибка возвращается в случае наличия фатальной ошибки из за которой продолжение работы не возможно
func (jbo *impl) Do() (err error) {
	var tp jobTypes.Type

	// Запрос конфигурации всех процессов
	if err = safeCall(jbo.getConfiguration); err != nil {
		return
	}
	// Составление списка процессов в соответствии с приоритетами запуска и остановки
	jbo.priority()
	// Запуск процессов с флагом Autostart в порядке приоритета
	for _, tp = range []jobTypes.Type{
		jobTypes.TypeForkWorker, jobTypes.TypeWorker, jobTypes.TypeTask} {
		if err = jbo.StartByType(tp); err != nil {
			return
		}
	}

	return
}

// Start Отправка команды запуска процесса
func (jbo *impl) Start(id string) (err error) {
	var prc *Process

	if _, prc, err = jbo.RegisteredProcessFindByID(id); err != nil {
		return
	}
	switch wrk := prc.P.(type) {
	case *jobTypes.Task:
		err = jbo.runTask(wrk)
	case *jobTypes.Worker:
		err = jbo.runWorker(wrk)
	case *jobTypes.ForkWorker:
		err = jbo.runForkWorker(wrk)
	default:
		err = Errors().TypeNotImplemented()
	}

	return
}

// StartByType Запуск процессов определённого типа в соответствии с очерёдностью запуска
func (jbo *impl) StartByType(tp jobTypes.Type) (err error) {
	var (
		prc   *Process
		st    *jobTypes.State
		isRun bool
		n     int
	)

	for n = range jbo.StartPriority {
		_, prc, err = jbo.RegisteredProcessFindByID(jbo.StartPriority[n])
		if err != nil && err != jbo.Errors().ProcessNotFound() {
			return
		}
		// Процесс не найден или не совпадает тип
		if err != nil || prc != nil && prc.Type != tp {
			continue
		}
		if st, err = prc.State(); err != nil {
			return
		}
		// Пропускаем всех с выключенным автостартом
		if !st.Conf.Autostart {
			continue
		}
		// Пропускаем уже запущенные
		if isRun, err = prc.IsRun(); err != nil {
			return
		}
		if isRun {
			continue
		}
		// Запуск процесса
		switch wrk := prc.P.(type) {
		case *jobTypes.Task:
			err = jbo.runTask(wrk)
		case *jobTypes.Worker:
			err = jbo.runWorker(wrk)
		case *jobTypes.ForkWorker:
			err = jbo.runForkWorker(wrk)
		default:
			err = jbo.Errors().TypeNotImplemented()
			return
		}
	}

	return
}

// IsCancelled Проверка состояния прерывания работы. Если передан не пустой id,
// тогда проверяется состояние для процесса, если передан пустой, то проверяется общее состояние для всех процессов.
// Истина - выполняется прерывание работы
// Ложь - разрешено нормальное выполнение процессов
func (jbo *impl) IsCancelled(id string) bool { return jbo.Exit.Load().(bool) }

// Cancel Сигнал завершения всех запущенных процессов
// Сигнал будет так же передан в подпроцессы запущенные как ForkWorker
func (jbo *impl) Cancel() { jbo.Event <- &jobEvent.Event{Act: jobEvent.ECancel} }
