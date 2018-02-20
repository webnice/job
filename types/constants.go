package types // import "gopkg.in/webnice/job.v1/types"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"time"
)

const (
	// LowPriopity Константа наименьшего приоритета
	LowPriopity = int32(^uint32(0) >> 1)

	// HighPriopity Константа наивысшего приоритета
	HighPriopity = int32(^1 << 30)

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
