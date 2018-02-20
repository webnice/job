package types // import "gopkg.in/webnice/job.v1/types"

//import "gopkg.in/webnice/debug.v1"
import log "gopkg.in/webnice/log.v2"
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

	log.Debug("New worker object")

	return jbo
}

// DestroyWorker Деструктор объектов Worker
func DestroyWorker(jbo *Worker) {

	log.Debug("Destroy worker object")

}
