package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"strings"

	"gopkg.in/webnice/job.v1/types"
)

// List Список зарегистрированных процессов
func (jbo *impl) List() (ret []Info) {
	var elm *list.Element
	var item *Process
	var ok bool

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*Process); !ok || item == nil {
			continue
		}
		switch wrk := item.P.(type) {
		case *types.Task:
			ret = append(ret, Info{
				ID:    wrk.ID,
				Type:  TTask,
				IsRun: wrk.State.IsRun.Load().(bool),
			})
		case *types.Worker:
			ret = append(ret, Info{
				ID:    wrk.ID,
				Type:  TWorker,
				IsRun: wrk.State.IsRun.Load().(bool),
			})
		case *types.ForkWorker:
			ret = append(ret, Info{
				ID:    wrk.ID,
				Type:  TForkWorker,
				IsRun: wrk.State.IsRun.Load().(bool),
			})
		}
	}

	return
}

// Find Поиск процесса по идентификатору
func (jbo *impl) Find(id string) (ret *Info, err error) {
	var all []Info
	var found bool
	var i int

	all = jbo.List()
	for i = range all {
		if strings.EqualFold(all[i].ID, id) {
			ret, found = &all[i], true
			break
		}
	}
	if !found {
		err = ErrorRegistredProcessNotFound()
	}

	return
}
