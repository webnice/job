package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
//	"gopkg.in/webnice/job.v1/types"
)

// ProcessDataExchange Функция вызывается из процесса для регистрации функции
// получения данных и получения функции отправки результата
// id - Идентификатор процесса, полученный процессом при вызове функции Info
// req - Функция получения данных, будет вызываться каждый раз, когда приходят внешние данные
// rsp - Функция отправки данных, можно вызывать многократно, для отправки данных
//func (jbo *impl) ProcessDataExchange(id string, req types.RequestFunc) (rsp types.ResponseFunc) {
//
//	return
//}
