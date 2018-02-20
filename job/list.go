package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"container/list"
	"strings"
)

// List Список зарегистрированных процессов
func (jbo *impl) List() (ret []Info) {
	var elm *list.Element
	var item *Process
	var ok bool

	for elm = jbo.ProcessList.Front(); elm != nil; elm = elm.Next() {
		if item, ok = elm.Value.(*Process); !ok {
			continue
		}
		switch {
		case item.Task != nil:
			ret = append(ret, Info{
				ID:    item.Task.ID,
				Type:  TTask,
				IsRun: item.Task.State.IsRun.Load().(bool),
			})
		case item.Worker != nil:
			ret = append(ret, Info{
				ID:    item.Worker.ID,
				Type:  TWorker,
				IsRun: item.Worker.State.IsRun.Load().(bool),
			})
		case item.ForkWorker != nil:
			ret = append(ret, Info{
				ID:    item.ForkWorker.ID,
				Type:  TForkWorker,
				IsRun: item.ForkWorker.State.IsRun.Load().(bool),
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
