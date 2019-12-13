package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"

	jobTypes "gopkg.in/webnice/job.v1/types"
)

// RegisterTask Регистрация простой управляемой задачи
func (jbo *impl) RegisterTask(obj Task) {
	var jb = jbo.Pool.TaskGet()
	jb.Self = obj
	jb.ID = getStructName(obj)
	jb.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: jb, Type: jobTypes.TypeTask})
}

// RegisterWorker Регистрация управляемого работника
func (jbo *impl) RegisterWorker(obj Worker) {
	var wk = jbo.Pool.WorkerGet()
	wk.Self = obj
	wk.ID = getStructName(obj)
	wk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: wk, Type: jobTypes.TypeWorker})
}

// RegisterForkWorker Регистрация управляемого работника
func (jbo *impl) RegisterForkWorker(obj ForkWorker) {
	var fwk = jbo.Pool.ForkWorkerGet()
	fwk.Self = obj
	fwk.ID = getStructName(obj)
	fwk.State.IsRun.Store(false)
	jbo.ProcessList.PushBack(&Process{P: fwk, Type: jobTypes.TypeForkWorker})
}

// Unregister Функция удаляет из реестра процессов процесс с указанным ID
// Для того чтобы быть удалённым, процесс должен быть в состоянии остановлен
func (jbo *impl) Unregister(id string) (err error) {
	var (
		item  *list.Element
		prc   *Process
		isRun bool
	)

	// Поиск зарегистрированного процесса по ID
	if item, prc, err = jbo.RegisteredProcessFindByID(id); err != nil {
		return
	}
	// Запущенный процесс нельзя разрегистрировать
	if isRun, err = prc.IsRun(); err != nil || isRun {
		if isRun {
			err = jbo.Errors().UnregisterProcessIsRunning()
		}
		return
	}
	// Удаление процесса из списка процессов
	_ = jbo.ProcessList.Remove(item)
	// Возврат объекта в пул
	if err = jbo.ProcessObjectReturnToPool(prc); err != nil {
		return
	}

	return
}
