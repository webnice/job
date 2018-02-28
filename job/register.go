package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"strings"

	"gopkg.in/webnice/job.v1/types"
)

// RegisterTask Регистрация простой управляемой задачи
func (jbo *impl) RegisterTask(obj Task) {
	var jb = jbo.Pool.TaskGet()
	jb.Self = obj
	jb.ID = getStructName(obj)
	jb.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: jb})
}

// RegisterWorker Регистрация управляемого работника
func (jbo *impl) RegisterWorker(obj Worker) {
	var wk = jbo.Pool.WorkerGet()
	wk.Self = obj
	wk.ID = getStructName(obj)
	wk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: wk})
}

// RegisterForkWorker Регистрация управляемого работника
func (jbo *impl) RegisterForkWorker(obj ForkWorker) {
	var fwk = jbo.Pool.ForkWorkerGet()
	fwk.Self = obj
	fwk.ID = getStructName(obj)
	fwk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: fwk})
}

// Unregister Функция удаляет из реестра процессов процесс с указанным ID
// Для того чтобы быть удалённым, процесс должен быть в состоянии остановлен
func (jbo *impl) Unregister(id string) (err error) {
	var elm, del *list.Element
	var prc *Process
	var ok, found, isRun bool
	var elmID string

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch wrk := prc.P.(type) {
		case *types.Task:
			elmID, isRun = wrk.ID, wrk.State.IsRun.Load().(bool)
		case *types.Worker:
			elmID, isRun = wrk.ID, wrk.State.IsRun.Load().(bool)
		case *types.ForkWorker:
			elmID, isRun = wrk.ID, wrk.State.IsRun.Load().(bool)
		}
		if strings.Index(elmID, id) == 0 {
			del, found = elm, true
			break
		}
	}
	if !found {
		err = ErrorRegistredProcessNotFound()
		return
	}
	if isRun {
		err = ErrorUnregisterProcessIsRunning()
		return
	}
	_ = jbo.ProcessList.Remove(del)
	// Возврат объекта в пул
	if prc, ok = del.Value.(*Process); ok && prc != nil {
		switch wrk := prc.P.(type) {
		case *types.Task:
			jbo.Pool.TaskPut(wrk)
		case *types.Worker:
			jbo.Pool.WorkerPut(wrk)
		case *types.ForkWorker:
			jbo.Pool.ForkWorkerPut(wrk)
		}
		prc.P = nil
	}

	return
}
