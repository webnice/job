package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"reflect"
	runtimeDebug "runtime/debug"
	"sync"
)

// Получение уникального имени пакета + имя структуры
func getStructName(obj interface{}) string {
	var rt reflect.Type
	var packageName, structureName string

	if rt = reflect.TypeOf(obj); rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	structureName = rt.Name()
	packageName = rt.PkgPath()

	return packageName + string(os.PathSeparator) + structureName
}

// Безопасный метод отправки пустой структуры (сигнала) в канал
func safeChannelSend(ch chan<- struct{}) {
	defer func() { _ = recover() }()
	ch <- struct{}{}
}

// Безопасно делает один вызов Done() для sync.WaitGroup
func safeWgDone(wg *sync.WaitGroup) {
	defer func() { _ = recover() }()
	wg.Done()
}

// Безопасно обнуляет до конца sync.WaitGroup
func safeWgDoneForAll(wg *sync.WaitGroup) {
	defer func() { _ = recover() }()
	for {
		wg.Done()
	}
}

// Безопасно выполняет Wait() для sync.WaitGroup
func safeWgWait(wg *sync.WaitGroup) {
	defer func() { _ = recover() }()
	wg.Wait()
}

// Безопасный запуск функции
func safeCall(fn func() error) (err error) {
	var ok bool

	defer func() {
		if e := recover(); e != nil {
			if err, ok = e.(error); !ok {
				err = fmt.Errorf("%v", e)
			}
			err = fmt.Errorf("%s\n%s", err, string(runtimeDebug.Stack()))
		}
	}()
	err = fn()

	return
}
