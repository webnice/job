package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	runtimeDebug "runtime/debug"
	"strings"
	"sync"
)

// Получение уникального имени пакета + имя структуры
func getStructName(obj interface{}) string {
	const callerSkip = 3
	var tmp []string
	var functionName, packageName, structureName string

	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		structureName = t.Elem().Name()
	} else {
		structureName = t.Name()
	}
	pc, _, _, ok := runtime.Caller(callerSkip)
	if ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			functionName = fn.Name()
		}
		tmp = strings.Split(functionName, string(os.PathSeparator))
		if len(tmp) > 1 {
			packageName += strings.Join(tmp[:len(tmp)-1], string(os.PathSeparator))
			functionName = tmp[len(tmp)-1]
		}
		tmp = strings.SplitN(functionName, `.`, 2)
		if len(tmp) == 2 {
			if packageName != "" {
				packageName += string(os.PathSeparator)
			}
			packageName += tmp[0]
			functionName = tmp[1]
		}
	}
	_ = functionName
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
