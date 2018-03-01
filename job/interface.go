package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"

	"gopkg.in/webnice/job.v1/types"
)

// Task Интерфейс простой управляемой задачи
type Task types.TaskInterface

// Worker Интерфейс управляемого работника
// Работник запускается в текущей копии приложения
type Worker types.WorkerInterface

// ForkWorker Интерфейс управляемого работника
// Работник запускается в новой копии приложения через syscall.ForkExec() - изолированный и убиваемый извне процесс
type ForkWorker types.ForkWorkerInterface

// Get Gets an interface of a package
func Get() Interface { return singleton }

// Error Последняя внутренняя ошибка
func Error() error { return singleton.Err }

// IsCancelled Проверка состояния прерывания работы
// Истина - выполняется прерывание работы всех воркеров
// Ложь - разрешено выполнение воркеров
func IsCancelled() bool { return singleton.IsCancelled() }

// Cancel Сигнал завершения всех запущенных процессов
func Cancel() { singleton.Cancel() }

// Do Запуск библиотеки, подготовка и запуск процессов
// Ошибка возвращается в случае наличия фатальной ошибки из за которой продолжение работы не возможно
func Do() error { return singleton.Do() }

// List of all registered processes
func List() []Info { return singleton.List() }

// RegisterErrorFunc Регистрация функции получения ошибок в работе управляемых процессов
func RegisterErrorFunc(fn OnErrorFunc) Interface { return singleton.RegisterErrorFunc(fn) }

// RegisterChangeStateFunc Регистрация функции получения изменения состояния процессов
func RegisterChangeStateFunc(fn OnChangeStateFunc) Interface {
	return singleton.RegisterChangeStateFunc(fn)
}

// UnregisterErrorFunc Удаление ранее зарегистрированной функции получения ошибок
func UnregisterErrorFunc() Interface { return singleton.UnregisterErrorFunc() }

// Unregister Функция удаляет из реестра процессов процесс с указанным ID
// Для того чтобы быть удалённым, процесс должен быть остановлен
func Unregister(id string) error { return singleton.Unregister(id) }

// Reset Сброс библиотеки, подготовка к повторному использованию
// Если были запущены процессы, то контроль над ними будет потерян
func Reset() { singleton.Reset() }

// Wait Ожидание завершения всех работающих процессов
func Wait() Interface { return singleton.Wait() }

// WaitWithTimeout Ожидание завершения всех работающих процессов, но не более чем время указанное в timeout
func WaitWithTimeout(timeout time.Duration) Interface { return singleton.WaitWithTimeout(timeout) }

// ProcessDataExchange Функция вызывается из процесса для регистрации функции
// получения данных и получения функции отправки результата
// id - Идентификатор процесса, полученный процессом при вызове функции Info
// req - Функция получения данных, будет вызываться каждый раз, когда приходят внешние данные
// rsp - Функция отправки данных, можно вызывать многократно, для отправки данных
//func ProcessDataExchange(id string, req types.RequestFunc) types.ResponseFunc {
//	return singleton.ProcessDataExchange(id, req)
//}
