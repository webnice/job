package pool

import (
	"sync"

	"github.com/webnice/job/types"
)

// New creates a new object and return interface.
func New() Interface {
	var poo = new(impl)
	poo.TaskPool = new(sync.Pool)
	poo.TaskPool.New = types.NewTask
	poo.WorkerPool = new(sync.Pool)
	poo.WorkerPool.New = types.NewWorker
	poo.ForkWorkerPool = new(sync.Pool)
	poo.ForkWorkerPool.New = types.NewForkWorker
	return poo
}

// TaskGet Получение из пула объектов, объекта управляющих данных, простой управляемой задачи.
func (poo *impl) TaskGet() *types.Task {
	return poo.TaskPool.Get().(*types.Task)
}

// TaskPut Возвращение в пул объектов, объекта управляющих данных, простой управляемой задачи.
func (poo *impl) TaskPut(jbo *types.Task) {
	poo.TaskPool.Put(jbo)
}

// WorkerGet Получение из пула объектов, объекта управляющих данных, управляемого работника.
func (poo *impl) WorkerGet() *types.Worker {
	return poo.WorkerPool.Get().(*types.Worker)
}

// WorkerPut Возвращение в пул объектов, объекта управляющих данных, управляемого работника.
func (poo *impl) WorkerPut(wrk *types.Worker) {
	poo.WorkerPool.Put(wrk)
}

// ForkWorkerGet Получение из пула объектов, объекта управляющих данных, управляемого работника.
func (poo *impl) ForkWorkerGet() *types.ForkWorker {
	return poo.ForkWorkerPool.Get().(*types.ForkWorker)
}

// ForkWorkerPut Возвращение в пул объектов, объекта управляющих данных, управляемого работника.
func (poo *impl) ForkWorkerPut(wrk *types.ForkWorker) {
	poo.ForkWorkerPool.Put(wrk)
}
