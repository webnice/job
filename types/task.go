package types

import (
	"context"
	"runtime"
)

// TaskInterface Интерфейс простой управляемой задачи
type TaskInterface interface {
	BaseInterface
}

// Task Структура управляющих данных простой управляемой задачи
type Task struct {
	Pith               // Общие для всех типов процессов переменные
	Self TaskInterface // Self
}

// NewTask Конструктор объектов Task
func NewTask() interface{} {
	var jbo = new(Task)
	jbo.Ctx, jbo.Cancel = context.WithCancel(context.Background())
	jbo.Self = nil
	runtime.SetFinalizer(jbo, DestroyTask)
	return jbo
}

// DestroyTask Деструктор объектов Task
func DestroyTask(jbo *Task) {
	//log.Debug("Destroy task object")
}
