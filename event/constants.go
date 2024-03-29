package event

const (
	// ECancel Событие остановки всех запущенных процессов.
	ECancel = Operation(`Cancel`)

	// ERestartProcess Событие перезапуска процесса завершившегося без ошибки.
	ERestartProcess = Operation(`Restart process`)

	// EProcessFatality Произошла фатальная ошибка, продолжение работы невозможно.
	EProcessFatality = Operation(`A fatal error has occurred, the continuation of the work is not possible`)

	// EOnError Событие ошибки.
	EOnError = Operation(`Process ends with error`)

	// ECancelError Вызов Cancel() завершился ошибкой.
	ECancelError = Operation(`Call cancel function ends with error`)

	// EProcessStarted Событие запуска процесса.
	EProcessStarted = Operation(`Process has started successfully`)

	// EProcessStopped Событие остановки процесса.
	EProcessStopped = Operation(`Process has stopped successfully`)

	// EStartProcess Событие запуска процесса.
	EStartProcess = Operation(`Start process`)
)
