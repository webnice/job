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
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case prc.Task != nil:
			prc.Task.State.Conf = prc.Task.Self.Info(prc.Task.ID)
			if prc.Task.State.Conf == nil {
				prc.Task.State.Conf = types.DefaultConfiguration()
			}
		case prc.Worker != nil:
			prc.Worker.State.Conf = prc.Worker.Self.Info(prc.Worker.ID)
			if prc.Worker.State.Conf == nil {
				prc.Worker.State.Conf = types.DefaultConfiguration()
			}
		case prc.ForkWorker != nil:
			prc.ForkWorker.State.Conf = prc.ForkWorker.Self.Info(prc.ForkWorker.ID)
			if prc.ForkWorker.State.Conf == nil {
				prc.ForkWorker.State.Conf = types.DefaultConfiguration()
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
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case prc.Task != nil:
			pst = append(pst, &plist{ID: prc.Task.ID, Start: prc.Task.State.Conf.PriorityStart, Stop: prc.Task.State.Conf.PriorityStop})
		case prc.Worker != nil:
			pst = append(pst, &plist{ID: prc.Worker.ID, Start: prc.Worker.State.Conf.PriorityStart, Stop: prc.Worker.State.Conf.PriorityStop})
		case prc.ForkWorker != nil:
			pst = append(pst, &plist{ID: prc.ForkWorker.ID, Start: prc.ForkWorker.State.Conf.PriorityStart, Stop: prc.ForkWorker.State.Conf.PriorityStop})
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
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case prc.Task != nil:
			if prc.Task.ID != id {
				continue
			}
			found = true
			err = jbo.runTask(prc.Task)
		case prc.Worker != nil:
			if prc.Worker.ID != id {
				continue
			}
			found = true
			err = jbo.runWorker(prc.Worker)
		case prc.ForkWorker != nil:
			if prc.ForkWorker.ID != id {
				continue
			}
			found = true
			err = jbo.runForkWorker(prc.ForkWorker)
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
			if prc, ok = elm.Value.(*Process); !ok {
				continue
			}
			if prc.ForkWorker == nil {
				continue
			}
			if id != "" && prc.ForkWorker.ID != id {
				continue
			} else if !prc.ForkWorker.State.Conf.Autostart ||
				prc.ForkWorker.ID != jbo.StartPriority[n] {
				continue
			}
			if err = jbo.runForkWorker(prc.ForkWorker); err != nil {
				return
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
			if prc, ok = elm.Value.(*Process); !ok {
				continue
			}
			if prc.Worker == nil {
				continue
			}
			if id != "" && prc.Worker.ID != id {
				continue
			} else if !prc.Worker.State.Conf.Autostart ||
				prc.Worker.ID != jbo.StartPriority[n] {
				continue
			}
			if err = jbo.runWorker(prc.Worker); err != nil {
				return
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
			if prc, ok = elm.Value.(*Process); !ok {
				continue
			}
			if prc.Task == nil {
				continue
			}
			if id != "" && prc.Task.ID != id {
				continue
			} else if !prc.Task.State.Conf.Autostart ||
				prc.Task.ID != jbo.StartPriority[n] {
				continue
			}
			if err = jbo.runTask(prc.Task); err != nil {
				return
			}
		}
	}

	return
}

// Cancel Сигнал завершения всех запущенных процессов
// Сигнал будет так же передан в подпроцессы запущенные как ForkWorker
func (jbo *impl) Cancel() { jbo.Event <- &event.Event{Act: event.ECancel} }
