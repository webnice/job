package types

import (
	"context"
	"runtime"
)

// WorkerInterface Интерфейс управляемого работника
// Работник запускается в текущей копии приложения
type WorkerInterface interface {
	BaseInterface

	// Prepare Функция выполнения действий подготавливающих воркер к работе
	// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
	Prepare() error
}

// Worker Структура управляющих данных управляемого работника
type Worker struct {
	Pith                 // Общие для всех типов процессов переменные
	Self WorkerInterface // Self
}

// NewWorker Конструктор объектов Worker
func NewWorker() interface{} {
	var jbo = new(Worker)
	jbo.Ctx, jbo.Cancel = context.WithCancel(context.Background())
	jbo.Self = nil
	runtime.SetFinalizer(jbo, DestroyWorker)
	return jbo
}

// DestroyWorker Деструктор объектов Worker
func DestroyWorker(jbo *Worker) {
	//log.Debug("Destroy worker object")
}
