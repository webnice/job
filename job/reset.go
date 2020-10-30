package job

import (
	"container/list"
	"context"
	"sync"

	jobEvent "github.com/webnice/job/event"
	jobPool "github.com/webnice/job/pool"
	jobTypes "github.com/webnice/job/types"
)

// Reset Сброс библиотеки, подготовка к повторному использованию
// Если были запущены процессы, то контроль над ними будет потерян
func (jbo *impl) Reset() Interface {
	// Контекст
	if jbo.CancelFunc != nil {
		jbo.CancelFunc()
	}
	jbo.TaskIDSync = new(sync.Mutex)
	jbo.Ctx, jbo.CancelFunc = context.WithCancel(context.Background())
	// Пул структур
	jbo.Pool = jobPool.New()
	// Список процессов
	switch jbo.ProcessList {
	case nil:
		jbo.ProcessList = list.New()
	default:
		jbo.ProcessList.Init()
	}
	// Группа ожидания
	if jbo.Wg != nil {
		safeWgDoneForAll(jbo.Wg)
	}
	jbo.Wg = new(sync.WaitGroup)
	// Канал событий
	if jbo.Event != nil {
		close(jbo.Event)
	}
	jbo.Event = make(chan *jobEvent.Event, jobTypes.EventBufLength)
	jbo.Exit.Store(false)

	return jbo
}
