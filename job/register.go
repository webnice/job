package job

import (
	"container/list"
	"fmt"
	"strconv"

	jobTypes "github.com/webnice/job/types"
)

// CreateTaskID Создание идентификатора задачи на базе объекта воркера.
func (jbo *impl) CreateTaskID(obj Task) (ret string) {
	var (
		err   error
		found []string
		max   uint64
		cur   uint64
		tmp   []string
		n     int
	)

	jbo.TaskIDSync.Lock()
	defer jbo.TaskIDSync.Unlock()
	ret = getStructName(obj)
	// Поиск всех совпадающих ID процессов.
	if err = jbo.RegisteredProcessIterate(
		func(elm *list.Element, prc *Process) (e error) {
			var (
				tid           string
				full, partial bool
			)
			if tid, e = prc.ID(); e != nil {
				return
			}
			if full, partial = jbo.compareID(tid, ret, 0); full || partial {
				found = append(found, tid)
			}
			return
		}); err != nil {
		return
	}
	if len(found) > 0 {
		for n = range found {
			if tmp = rexNameMatch.FindStringSubmatch(found[n]); len(tmp) != 4 {
				continue
			}
			cur, _ = strconv.ParseUint(tmp[3], 10, 64)
			if cur > max {
				max = cur
			}
		}
		max++
		ret = fmt.Sprintf("%s-%d", ret, max)
	}

	return
}

// RegisterTask Регистрация простой управляемой задачи.
func (jbo *impl) RegisterTask(obj Task) (ret string) {
	var jb = jbo.Pool.TaskGet()
	jb.Self = obj
	jb.ID = jbo.CreateTaskID(obj)
	jb.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: jb, Type: jobTypes.TypeTask})
	return jb.ID
}

// RegisterWorker Регистрация управляемого работника.
func (jbo *impl) RegisterWorker(obj Worker) (ret string) {
	var wk = jbo.Pool.WorkerGet()
	wk.Self = obj
	wk.ID = jbo.CreateTaskID(obj)
	wk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: wk, Type: jobTypes.TypeWorker})
	return wk.ID
}

// RegisterForkWorker Регистрация управляемого работника.
func (jbo *impl) RegisterForkWorker(obj ForkWorker) (ret string) {
	var fwk = jbo.Pool.ForkWorkerGet()
	fwk.Self = obj
	fwk.ID = jbo.CreateTaskID(obj)
	fwk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: fwk, Type: jobTypes.TypeForkWorker})
	return fwk.ID
}

// Unregister Функция удаляет из реестра процессов процесс с указанным ID.
// Для того чтобы быть удалённым, процесс должен быть в состоянии остановлен.
func (jbo *impl) Unregister(id string) (err error) {
	var (
		item  *list.Element
		prc   *Process
		isRun bool
	)

	// Поиск зарегистрированного процесса по ID.
	if item, prc, err = jbo.RegisteredProcessFindByID(id); err != nil {
		return
	}
	// Запущенный процесс нельзя разрегистрировать.
	if isRun, err = prc.IsRun(); err != nil || isRun {
		if isRun {
			err = jbo.Errors().UnregisterProcessIsRunning()
		}
		return
	}
	// Удаление процесса из списка процессов.
	_ = jbo.ProcessList.Remove(item)
	// Возврат объекта в пул.
	if err = jbo.ProcessObjectReturnToPool(prc); err != nil {
		return
	}

	return
}
