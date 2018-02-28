package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"sort"

	"gopkg.in/webnice/job.v1/event"
	"gopkg.in/webnice/job.v1/types"
)

// Запрос конфигурации всех процессов
func (jbo *impl) getConfiguration() (err error) {
	var elm *list.Element
	var ok bool
	var prc *Process

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok || prc == nil {
			continue
		}
		switch wrk := prc.P.(type) {
		case *types.Task:
			wrk.State.Conf = wrk.Self.Info(wrk.ID)
			if wrk.State.Conf == nil {
				wrk.State.Conf = types.DefaultConfiguration()
			}
		case *types.Worker:
			wrk.State.Conf = wrk.Self.Info(wrk.ID)
			if wrk.State.Conf == nil {
				wrk.State.Conf = types.DefaultConfiguration()
			}
		case *types.ForkWorker:
			wrk.State.Conf = wrk.Self.Info(wrk.ID)
			if wrk.State.Conf == nil {
				wrk.State.Conf = types.DefaultConfiguration()
			}
		}
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
	var elm *list.Element
	var pst []*plist
	var ok bool
	var prc *Process
	var n int

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok || prc == nil {
			continue
		}
		switch wrk := prc.P.(type) {
		case *types.Task:
			pst = append(pst, &plist{ID: wrk.ID, Start: wrk.State.Conf.PriorityStart, Stop: wrk.State.Conf.PriorityStop})
		case *types.Worker:
			pst = append(pst, &plist{ID: wrk.ID, Start: wrk.State.Conf.PriorityStart, Stop: wrk.State.Conf.PriorityStop})
		case *types.ForkWorker:
			pst = append(pst, &plist{ID: wrk.ID, Start: wrk.State.Conf.PriorityStart, Stop: wrk.State.Conf.PriorityStop})
		}
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

// Do Запуск библиотеки, подготовка и запуск процессов
// Ошибка возвращается в случае наличия фатальной ошибки из за которой продолжение работы не возможно
func (jbo *impl) Do() (err error) {
	// Запрос конфигурации всех процессов
	if err = safeCall(jbo.getConfiguration); err != nil {
		return
	}
	// Составление списка процессов в соответствии с приоритетами запуска и остановки
	jbo.priority()
	// Запуск процессов frokworker
	if err = jbo.doFrokWorker(""); err != nil {
		return
	}
	// Запуск процессов worker
	if err = jbo.doWorker(""); err != nil {
		return
	}
	// Запуск процессов task
	if err = jbo.doTask(""); err != nil {
		return
	}

	return
}

// Start Отправка команды запуска процесса
func (jbo *impl) Start(id string) (err error) {
	var elm *list.Element
	var prc *Process
	var ok, found bool

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok || prc == nil {
			continue
		}
		switch wrk := prc.P.(type) {
		case *types.Task:
			if wrk.ID != id {
				continue
			}
			err, found = jbo.runTask(wrk), true
		case *types.Worker:
			if wrk.ID != id {
				continue
			}
			err, found = jbo.runWorker(wrk), true
		case *types.ForkWorker:
			if wrk.ID != id {
				continue
			}
			err, found = jbo.runForkWorker(wrk), true
		}
	}
	if !found {
		err = ErrorProcessNotFound()
		return
	}

	return
}

// Запуск всех forkworker процессов
// Если id="" - Запускаются все процессы с флагом Autostart=true
// Если id указан, запускается процесс с указанным id
func (jbo *impl) doFrokWorker(id string) (err error) {
	var elm *list.Element
	var prc *Process
	var n int
	var ok bool

	for n = range jbo.StartPriority {
		for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
			if prc, ok = elm.Value.(*Process); !ok || prc == nil {
				continue
			}
			switch wrk := prc.P.(type) {
			case *types.ForkWorker:
				if id != "" && wrk.ID != id {
					continue
				} else if !wrk.State.Conf.Autostart ||
					wrk.ID != jbo.StartPriority[n] {
					continue
				}
				if err = jbo.runForkWorker(wrk); err != nil {
					return
				}

			}
		}
	}

	return
}

// Запуск всех worker процессов
// Если id="" - Запускаются все процессы с флагом Autostart=true
// Если id указан, запускается процесс с указанным id
func (jbo *impl) doWorker(id string) (err error) {
	var elm *list.Element
	var prc *Process
	var n int
	var ok bool

	for n = range jbo.StartPriority {
		for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
			if prc, ok = elm.Value.(*Process); !ok || prc == nil {
				continue
			}
			switch wrk := prc.P.(type) {
			case *types.Worker:
				if id != "" && wrk.ID != id {
					continue
				} else if !wrk.State.Conf.Autostart ||
					wrk.ID != jbo.StartPriority[n] {
					continue
				}
				if err = jbo.runWorker(wrk); err != nil {
					return
				}
			}
		}
	}

	return
}

// Запуск task процессов
// Если id="" - Запускаются все процессы с флагом Autostart=true
// Если id указан, запускается процесс с указанным id
func (jbo *impl) doTask(id string) (err error) {
	var elm *list.Element
	var prc *Process
	var n int
	var ok bool

	for n = range jbo.StartPriority {
		for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
			if prc, ok = elm.Value.(*Process); !ok || prc == nil {
				continue
			}
			switch wrk := prc.P.(type) {
			case *types.Task:
				if id != "" && wrk.ID != id {
					continue
				} else if !wrk.State.Conf.Autostart ||
					wrk.ID != jbo.StartPriority[n] {
					continue
				}
				if err = jbo.runTask(wrk); err != nil {
					return
				}
			}
		}
	}

	return
}

// Cancel Сигнал завершения всех запущенных процессов
// Сигнал будет так же передан в подпроцессы запущенные как ForkWorker
func (jbo *impl) Cancel() { jbo.Event <- &event.Event{Act: event.ECancel} }
