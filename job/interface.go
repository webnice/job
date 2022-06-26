package job

import (
	"time"

	jobTypes "github.com/webnice/job/types"
)

// Task Интерфейс простой управляемой задачи.
type Task jobTypes.TaskInterface

// Worker Интерфейс управляемого работника.
// Работник запускается в текущей копии приложения.
type Worker jobTypes.WorkerInterface

// ForkWorker Интерфейс управляемого работника.
// Работник запускается в новой копии приложения через syscall.ForkExec() - изолированный и убиваемый извне процесс.
type ForkWorker jobTypes.ForkWorkerInterface

// Get an interface of a package.
func Get() Interface { return singleton }

// Err Последняя внутренняя ошибка.
func Err() error { return singleton.err }

// IsCancelled Проверка состояния прерывания работы. Если передан не пустой id, то проверяется состояние для процесса,
// если передан пустой, то проверяется общее состояние для всех процессов.
// Истина - выполняется прерывание работы.
// Ложь - разрешено нормальное выполнение процессов.
func IsCancelled(id string) bool { return singleton.IsCancelled(id) }

// Cancel Сигнал завершения всех запущенных процессов.
func Cancel() { singleton.Cancel() }

// Do Запуск библиотеки, подготовка и запуск процессов.
// Ошибка возвращается в случае наличия фатальной ошибки, из-за которой продолжение работы невозможно.
func Do() error { return singleton.Do() }

// List of all registered processes
func List() []*Info { return singleton.List() }

// RegisterErrorFunc Регистрация функции получения ошибок в работе управляемых процессов.
func RegisterErrorFunc(fn OnErrorFunc) Interface { return singleton.RegisterErrorFunc(fn) }

// RegisterChangeStateFunc Регистрация функции получения изменения состояния процессов.
func RegisterChangeStateFunc(fn OnChangeStateFunc) Interface {
	return singleton.RegisterChangeStateFunc(fn)
}

// UnregisterErrorFunc Удаление ранее зарегистрированной функции получения ошибок.
func UnregisterErrorFunc() Interface { return singleton.UnregisterErrorFunc() }

// Unregister Функция удаляет из реестра процессов процесс с указанным ID.
// Для того чтобы быть удалённым, процесс должен быть остановлен.
func Unregister(id string) error { return singleton.Unregister(id) }

// Reset Сброс библиотеки, подготовка к повторному использованию.
// Если были запущены процессы, то контроль над ними будет потерян.
func Reset() { singleton.Reset() }

// Wait Ожидание завершения всех работающих процессов.
func Wait() Interface { return singleton.Wait() }

// WaitWithTimeout Ожидание завершения всех работающих процессов, но не более чем время указанное в timeout.
func WaitWithTimeout(timeout time.Duration) Interface { return singleton.WaitWithTimeout(timeout) }
