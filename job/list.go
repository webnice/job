package job

import (
	"container/list"
	"fmt"
	"strings"
)

// List Список зарегистрированных процессов
func (jbo *impl) List() (ret []*Info) {
	ret = make([]*Info, 0, jbo.ProcessList.Len())
	jbo.err = jbo.RegisteredProcessIterate(func(elm *list.Element, prc *Process) (e error) {
		var item = new(Info)

		item.Type = prc.Type
		if item.ID, e = prc.ID(); e != nil {
			return
		}
		if item.IsRun, e = prc.IsRun(); e != nil {
			return
		}
		ret = append(ret, item)
		return
	})

	return
}

// Сравнение идентификатора с искомым идентификатором и расширением
// Процессы именуются по имени пакета и имени типа объекта в пакете,
// если существует не один объект процесса, то все остальные дополняются числом через дефис:
// - application/workers/web/impl
// - application/workers/web/impl-1
// - application/workers/web/impl-2
// При сравнении указывается в id сравниваемый и в search искомый идентификатор.
// В ext указывается искомое значение дополняемого числа
//
// Функция вернёт
// - full=true                  - если id полностью совпадает с search
// - full=false и partial=true  - если id совпадает с search, но дополнение ext не совпадает с указанным в id
// - full=false и partial=false - при полном не совпадении
func (jbo *impl) compareID(id string, search string, ext uint64) (full bool, partial bool) {
	var (
		searchFull string
		tmp        []string
	)

	if searchFull = search; ext > 0 {
		searchFull = fmt.Sprintf("%s-%d", search, ext)
	}
	if strings.EqualFold(id, searchFull) {
		full = true
		return
	}
	if tmp = rexNameMatch.FindStringSubmatch(id); len(tmp) != 4 {
		return
	}
	if strings.EqualFold(tmp[1], search) {
		partial = true
	}
	if strings.EqualFold(tmp[3], fmt.Sprintf("%d", ext)) {
		full = true
	}

	return
}

// Find Поиск процесса по идентификатору
func (jbo *impl) Find(id string) (ret *Info, err error) {
	var (
		all   []*Info
		found bool
		i     int
	)

	all = jbo.List()
	for i = range all {
		if f, p := jbo.compareID(all[i].ID, id, 0); f || p {
			ret, found = all[i], true
			break
		}
	}
	if !found || ret == nil {
		err = jbo.Errors().RegisteredProcessNotFound()
	}

	return
}
