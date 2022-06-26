package pool

import (
	"sync"

	"github.com/webnice/job/types"
)

// Interface is an interface of package.
type Interface interface {
	// TaskGet Получение из пула объектов, объекта управляющих данных, простой управляемой задачи.
	TaskGet() *types.Task

	// TaskPut Возвращение в пул объектов, объекта управляющих данных, простой управляемой задачи.
	TaskPut(jbo *types.Task)

	// WorkerGet Получение из пула объектов, объекта управляющих данных, управляемого работника.
	WorkerGet() *types.Worker

	// WorkerPut Возвращение в пул объектов, объекта управляющих данных, управляемого работника.
	WorkerPut(wrk *types.Worker)

	// ForkWorkerGet Получение из пула объектов, объекта управляющих данных, управляемого работника.
	ForkWorkerGet() *types.ForkWorker

	// ForkWorkerPut Возвращение в пул объектов, объекта управляющих данных, управляемого работника.
	ForkWorkerPut(wrk *types.ForkWorker)
}

// impl is an implementation of package.
type impl struct {
	TaskPool       *sync.Pool // Объекты *types.Task
	WorkerPool     *sync.Pool // Объекты *types.Worker
	ForkWorkerPool *sync.Pool // Объекты *types.ForkWorker
}
