package types // import "gopkg.in/webnice/job.v1/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"context"
	"sync/atomic"
	"time"
)

// BaseInterface Базовый интерфейс всех воркеров
type BaseInterface interface {
	// Info Функция конфигурации процесса,
	// - в функцию передаётся уникальный идентификатор присвоенный процессу
	// - функция должна вернуть конфигурацию или nil, если возвращается nil, то применяется конфигурация по умолчанию
	Info(id string) *Configuration

	// Cancel Функция прерывания работы
	Cancel() error

	// Worker Функция-реализация процесса, данная функция будет запущена в горутине
	// до тех пор пока функция не завершился воркер считается работающим
	Worker() error
}

// RequestFunc Функция получения данных из процесса
type RequestFunc func(interface{})

// ResponseFunc Функция передачи данных результата работы процеса
type ResponseFunc func(interface{})

// Configuration Конфигурация процессов
type Configuration struct {
	// Autostart Флаг автоматического запуска процесса по готовности приложения
	// - если =true - процесс запускается автоматически
	// - если =false - процесс запускается только через отправку сигнала запуска процесса
	Autostart bool

	// Restart Флаг перезапуска приложения
	// - если =true - процесс будет перезапущен автоматически если завершится без ошибки
	// - если =false - остановившийся процесс остаётся завершенным
	Restart bool

	// Fatality Флаг фатальности. Если процесс завершится с ошибкой, то будет вызван сигнал остановки приложения
	// - если =true - в случае остановки процесса с ошибкой, сигнал остановки отправляется всем другим процессам приложения
	// - если =false - остановившийся процесс остаётся завершенным вне зависимости завершился он с ошибкой или без
	Fatality bool

	// PriorityStart Приоритет запуска, от меньшего к большему
	PriorityStart int32

	// PriorityStop Приоритет остановки, от меньшего к большему
	PriorityStop int32

	// RestartTimeout Таймаут между перезапусками процесса, для конфигурации в которой указан AutoRestart=true
	RestartTimeout time.Duration

	// KillTimeout Таймаут выполнения команды kill, если выполнение функции Cancel() не привело к завершению процесса
	// Используется только для ForkWorker
	KillTimeout time.Duration

	// CancelTimeout Таймаут завершения процесса после выполнения функции Cancel()
	// Если больше 0, то процесс считается завершенным, несмотря на то что он всё еще работает
	// ВНИМАНИЕ, Важно понимать когда можно устанавливать это значение, так как CancelTimeout может привести к утечке горутин
	CancelTimeout time.Duration
}

// Pith Основная структура данных процеса
type Pith struct {
	Ctx    context.Context    // Контекст
	Cancel context.CancelFunc // Функция прерывания выполнения
	ID     string             // пакет + название структуры = уникальный ID воркера
	State  State              // Состояние процесса
}

// State of process
type State struct {
	IsRun atomic.Value   // (bool) =true - процесс запущен, =false - процесс остановлен
	Conf  *Configuration // Конфигурация процесса
}
