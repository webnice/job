package job

import (
	"container/list"
	"context"
	"regexp"
	"sync"
	"sync/atomic"
	"time"

	jobEvent "github.com/webnice/job/v2/event"
	jobPool "github.com/webnice/job/v2/pool"
	jobTypes "github.com/webnice/job/v2/types"
)

var (
	singleton *impl

	// Регэксп для сравнения имени процесса.
	rexNameMatch = regexp.MustCompile(`(?ms)^(.+?)(\-(\d+))*$`)
)

// Interface is an interface of job package
type Interface interface {
	// Do Запуск библиотеки, подготовка и запуск процессов с флагом Autostart.
	// Ошибка возвращается в случае наличия фатальной ошибки, из-за которой продолжение работы невозможно.
	Do() error

	// Start Отправка команды запуска процесса
	Start(id string) error

	// IsCancelled Проверка состояния прерывания работы. Если передан не пустой id,
	// тогда проверяется состояние для процесса, если передан пустой, то проверяется общее состояние для всех процессов.
	// Истина - выполняется прерывание работы.
	// Ложь - разрешено нормальное выполнение процессов.
	IsCancelled(id string) bool

	// Cancel Сигнал завершения всех запущенных процессов.
	Cancel()

	// List Список зарегистрированных процессов.
	List() []*Info

	// RegisterErrorFunc Регистрация функции получения ошибок о работе управляемых процессов.
	RegisterErrorFunc(fn OnErrorFunc) Interface

	// RegisterChangeStateFunc Регистрация функции получения изменения состояния процессов.
	RegisterChangeStateFunc(fn OnChangeStateFunc) Interface

	// RegisterTask Регистрация простой управляемой задачи.
	RegisterTask(obj Task) (ret string)

	// RegisterWorker Регистрация управляемого работника.
	RegisterWorker(obj Worker) (ret string)

	// RegisterForkWorker Регистрация управляемого работника.
	RegisterForkWorker(obj ForkWorker) (ret string)

	// UnregisterErrorFunc Удаление ранее зарегистрированной функции получения ошибок.
	UnregisterErrorFunc() Interface

	// Unregister Функция удаляет из реестра процессов процесс с указанным ID.
	// Для того чтобы быть удалённым, процесс должен быть в состоянии остановлен.
	Unregister(id string) error

	// Reset Сброс библиотеки, подготовка к повторному использованию.
	// Если были запущены процессы, то контроль над ними будет потерян.
	Reset() Interface

	// Wait Ожидание завершения всех работающих процессов.
	Wait() Interface

	// WaitWithTimeout Ожидание завершения всех работающих процессов, но не более чем время указанное в timeout.
	WaitWithTimeout(timeout time.Duration) Interface

	// ОШИБКИ

	// Err Последняя внутренняя ошибка.
	Err() error

	// Errors Все ошибки известного состояния, которые могут вернуть функции пакета.
	Errors() *Error
}

// impl is an implementation of package.
type impl struct {
	err             error                // Последняя внутренняя ошибка.
	Ctx             context.Context      // Context of package.
	CancelFunc      context.CancelFunc   // Функция прерывания ожидания.
	Pool            jobPool.Interface    // Интерфейс пула управляющих данных.
	ProcessList     *list.List           // List of jobs, хранит объекты *Process.
	ErrorFunc       OnErrorFunc          // Функция получения ошибок процессов.
	ChangeStateFunc OnChangeStateFunc    // Функция изменения состояния процессов.
	Wg              *sync.WaitGroup      // Ожидание завершения всех выполняющихся задач.
	StartPriority   []string             // Списко идентификаторов процессов отсортированных в порядке приоритета запуска.
	StopPriority    []string             // Списко идентификаторов процессов отсортированных в порядке приоритета остановки.
	Event           chan *jobEvent.Event // Канал внутренних событий.
	Exit            atomic.Value         // Если =true, выполняется завершение работы. Запрещён запуск и перезапуск процессов.
	TaskIDSync      *sync.Mutex          // Мьютекс для генерации ID задачи.
}

// Process Ссылка на процесс.
type Process struct {
	P    interface{}   // Присваиваются: *types.Task, *types.Worker, *types.ForkWorker.
	Type jobTypes.Type // Тип процесса, значения: task, worker, forkworker.
}

// OnErrorFunc Функция получения ошибок.
type OnErrorFunc func(id string, err error)

// OnChangeStateFunc Функция изменения состояния процессов.
type OnChangeStateFunc func(id string, running bool)

// Info information about workers in memory.
type Info struct {
	ID    string        // Идентификатор процесса.
	Type  jobTypes.Type // Тип процесса, значения: task, worker, forkworker.
	IsRun bool          // =true - процесс запущен.
}

// ID Идентификатор процесса.
type ID struct {
	name         string
	serialNumber uint64
	pid          int64
}
