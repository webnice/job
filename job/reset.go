package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"context"
	"sync"

	"gopkg.in/webnice/job.v1/event"
	"gopkg.in/webnice/job.v1/pool"
)

// Reset Сброс библиотеки, подготовка к повторному использованию
// Если были запущены процессы, то контроль над ними будет потерян
func (jbo *impl) Reset() {
	// Контекст
	if jbo.CancelFunc != nil {
		jbo.CancelFunc()
	}
	jbo.Ctx, jbo.CancelFunc = context.WithCancel(context.Background())
	// Пул структур
	jbo.Pool = pool.New()
	// Список процессов
	if jbo.ProcessList == nil {
		jbo.ProcessList = list.New()
	} else {
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
	jbo.Event = make(chan *event.Event, _EventBufLength)
	jbo.Exit.Store(false)
}
