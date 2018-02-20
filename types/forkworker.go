package types // import "gopkg.in/webnice/job.v1/types"

//import debug "gopkg.in/webnice/debug.v1"
import log "gopkg.in/webnice/log.v2"
import (
	"context"
	"runtime"
)

// ForkWorkerInterface Интерфейс управляемого работника
// Работник запускается в новой копии приложения - изолированный процесс
type ForkWorkerInterface interface {
	BaseInterface

	// Prepare Функция выполнения действий подготавливающих воркер к работе
	// Завершение с ошибкой означает, что процесс не удалось подготовить к запуску
	Prepare() error

	// Interrupt Функция вызывается при получении извне сигнала прерывания работы INT (interrupt) (2)
	Interrupt() error
}

// ForkWorker Структура управляющих данных управляемого работника
type ForkWorker struct {
	Pith                     // Общие для всех типов процессов переменные
	Self ForkWorkerInterface // Self
}

// NewForkWorker Конструктор объектов ForkWorker
func NewForkWorker() interface{} {
	var jbo = new(ForkWorker)
	jbo.Ctx, jbo.Cancel = context.WithCancel(context.Background())
	jbo.Self = nil
	runtime.SetFinalizer(jbo, DestroyForkWorker)

	log.Debug("New fork worker object")

	return jbo
}

// DestroyForkWorker Деструктор объектов Worker
func DestroyForkWorker(jbo *ForkWorker) {

	log.Debug("Destroy fork worker object")

}
