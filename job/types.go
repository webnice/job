package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"context"
	"sync"
	"sync/atomic"
	"time"

	"gopkg.in/webnice/job.v1/event"
	"gopkg.in/webnice/job.v1/pool"
)

const (
	// TTask Процесс является task
	TTask = Type(`task`)

	// TWorker Процесс является worker
	TWorker = Type(`worker`)

	// TForkWorker Процесс является forkworker
	TForkWorker = Type(`forkworker`)

	_EventBufLength = int(10000)
)

var singleton *impl

// Interface is an interface of job package
type Interface interface {
	// Do Запуск библиотеки, подготовка и запуск процессов
	// Ошибка возвращается в случае наличия фатальной ошибки из за которой продолжение работы не возможно
	Do() error

	// Start Отправка команды запуска процесса
	Start(id string) error

	// Error Последняя внутренняя ошибка
	Error() error

	// Cancel Сигнал завершения всех запущенных процессов
	Cancel()

	// List Список зарегистрированных процессов
	List() []Info

	// RegisterErrorFunc Регистрация функции получения ошибок в работе управляемых процессов
	RegisterErrorFunc(fn OnErrorFunc) Interface

	// RegisterChangeStateFunc Регистрация функции получения изменения состояния процессов
	RegisterChangeStateFunc(fn OnChangeStateFunc) Interface

	// RegisterTask Регистрация простой управляемой задачи
	RegisterTask(obj Task)

	// RegisterWorker Регистрация управляемого работника
	RegisterWorker(obj Worker)

	// RegisterForkWorker Регистрация управляемого работника
	RegisterForkWorker(obj ForkWorker)

	// UnregisterErrorFunc Удаление ранее зарегистрированной функции получения ошибок
	UnregisterErrorFunc() Interface

	// Unregister Функция удаляет из реестра процессов процесс с указанным ID
	// Для того чтобы быть удалённым, процесс должен быть в состоянии остановлен
	Unregister(id string) error

	// Reset Сброс библиотеки, подготовка к повторному использованию
	// Если были запущены процессы, то контроль над ними будет потерян
	Reset()

	// Wait Ожидание завершения всех работающих процессов
	Wait() Interface

	// WaitWithTimeout Ожидание завершения всех работающих процессов, но не более чем время указанное в timeout
	WaitWithTimeout(timeout time.Duration) Interface

	// ProcessDataExchange Функция вызывается из процесса для регистрации функции
	// получения данных и получения функции отправки результата
	// id - Идентификатор процесса, полученный процессом при вызове функции Info
	// req - Функция получения данных, будет вызываться каждый раз, когда приходят внешние данные
	// rsp - Функция отправки данных, можно вызывать многократно, для отправки данных
	//ProcessDataExchange(id string, req types.RequestFunc) types.ResponseFunc
}

// impl is an implementation of package
type impl struct {
	Ctx             context.Context    // Context of package
	CancelFunc      context.CancelFunc // Функция прерывания ожидания
	Pool            pool.Interface     // Интерфейс пула управляющих данных
	ProcessList     *list.List         // List of jobs, хранит объекты *Process
	ErrorFunc       OnErrorFunc        // Функция получения ошибок процессов
	ChangeStateFunc OnChangeStateFunc  // Функция изменения состояния процессов
	Err             error              // Последняя внутренняя ошибка
	Wg              *sync.WaitGroup    // Ожидание завершения всех выполняющихся задач
	StartPriority   []string           // Списко идентификаторов процессов отсортированных в порядке приоритета запуска
	StopPriority    []string           // Списко идентификаторов процессов отсортированных в порядке приоритета остановки
	Event           chan *event.Event  // Канал внутренних событий
	Exit            atomic.Value       // Если =true, выполняется завершение работы. Запрещён запуск и перезапуск процессов
}

// Process Ссылка на процесс
type Process struct {
	// P Присваиваются: *types.Task, *types.Worker, *types.ForkWorker
	P interface{}
}

// OnErrorFunc Функция получения ошибок
type OnErrorFunc func(id string, err error)

// OnChangeStateFunc Функция изменения состояния процессов
type OnChangeStateFunc func(id string, running bool)

// Info information about workers in memory
type Info struct {
	ID    string // Идентификатор процесса
	Type  Type   // Тип процесса, значения: job, worker, forkworker
	IsRun bool   // =true - процесс запущен
}

// Type Тип процесса
type Type string

// String Convert type to string
func (tpe Type) String() string { return string(tpe) }
