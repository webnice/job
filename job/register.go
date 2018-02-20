package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"strings"
)

// RegisterTask Регистрация простой управляемой задачи
func (jbo *impl) RegisterTask(obj Task) {
	var jb = jbo.Pool.TaskGet()
	jb.Self = obj
	jb.ID = getStructName(obj)
	jb.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{Task: jb})
}

// RegisterWorker Регистрация управляемого работника
func (jbo *impl) RegisterWorker(obj Worker) {
	var wk = jbo.Pool.WorkerGet()
	wk.Self = obj
	wk.ID = getStructName(obj)
	wk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{Worker: wk})
}

// RegisterForkWorker Регистрация управляемого работника
func (jbo *impl) RegisterForkWorker(obj ForkWorker) {
	var fwk = jbo.Pool.ForkWorkerGet()
	fwk.Self = obj
	fwk.ID = getStructName(obj)
	fwk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{ForkWorker: fwk})
}

// Unregister Функция удаляет из реестра процессов процесс с указанным ID
// Для того чтобы быть удалённым, процесс должен быть в состоянии остановлен
func (jbo *impl) Unregister(id string) (err error) {
	var elm, del *list.Element
	var prc *Process
	var ok, found, isRun, t, w, f bool
	var elmID string

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if prc, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case prc.Task != nil:
			elmID, isRun, t = prc.Task.ID, prc.Task.State.IsRun.Load().(bool), true
		case prc.Worker != nil:
			elmID, isRun, w = prc.Worker.ID, prc.Worker.State.IsRun.Load().(bool), true
		case prc.ForkWorker != nil:
			elmID, isRun, f = prc.ForkWorker.ID, prc.ForkWorker.State.IsRun.Load().(bool), true
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
	if prc, ok = del.Value.(*Process); ok {
		switch {
		case t:
			jbo.Pool.TaskPut(prc.Task)
			prc.Task = nil
		case w:
			jbo.Pool.WorkerPut(prc.Worker)
			prc.Worker = nil
		case f:
			jbo.Pool.ForkWorkerPut(prc.ForkWorker)
			prc.ForkWorker = nil
		}
	}

	return
}
