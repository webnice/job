package types // import "gopkg.in/webnice/job.v1/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "time"

const (
	// TypeTask Процесс является task
	TypeTask = Type(`task`)

	// TypeWorker Процесс является worker
	TypeWorker = Type(`worker`)

	// TypeForkWorker Процесс является forkworker
	TypeForkWorker = Type(`forkworker`)

	// Размер буфера событий
	EventBufLength = int(10000)

	// LowPriority Константа наименьшего приоритета
	LowPriority = int32(^uint32(0) >> 1)

	// HighPriority Константа наивысшего приоритета
	HighPriority = int32(^1 << 30)

	// DefaultRestartTimeout Время ожидания между перезапусками остановившегося процесса
	DefaultRestartTimeout = time.Second / 4

	// DefaultKillTimeout Время ожидания остановки процесса ForkWorker перед отправкой не игнорируемого сигнала KILL процессу
	DefaultKillTimeout = time.Second * 15
)

// Конфигурация по умолчанию для всех процессов не вернувших свою конфигурацию
var defaultConfiguration = &Configuration{
	Autostart:      true,
	RestartTimeout: DefaultRestartTimeout,
	KillTimeout:    DefaultKillTimeout,
}

// Type Тип процесса
type Type string

// String Convert type to string
func (tpe Type) String() string { return string(tpe) }
