package event // import "gopkg.in/webnice/job.v1/event"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"

const (
	// ECancel Событие остановки всех запущенных процессов
	ECancel = Operation(`Cancel`)

	// ERestartProcess Событие перезапуска процесса завершившегося без ошибки
	ERestartProcess = Operation(`Restart process`)

	// EProcessFatality Произошла фатальная ошибка, продолжение работы не возможно
	EProcessFatality = Operation(`A fatal error has occurred, the continuation of the work is not possible`)

	// EOnError Событие ошибки
	EOnError = Operation(`Process ends with error`)

	// ECancelError Вызов Cancel() завершился ошибкой
	ECancelError = Operation(`Call cancel function ends with error`)

	// EProcessStarted Событие запуска процесса
	EProcessStarted = Operation(`Process has started successfully`)

	// EProcessStoped Событие остановки процесса
	EProcessStoped = Operation(`Process has stopped successfully`)

	// EStartProcess Событие запуска процесса
	EStartProcess = Operation(`Start process`)
)
