package job

// Все ошибки определены как константы.
const (
	cUnexpectedError            = "Unexpected error"
	cNotImplemented             = "Not implemented"
	cTypeNotImplemented         = "Type of process is not implemented"
	cUnregisterProcessIsRunning = "The process is current running, you must first stop the process"
	cRegisteredProcessNotFound  = "Registered process with the specified identifier was not found"
	cDeadlineExceeded           = "Deadline exceeded"
	cProcessesAreStillRunning   = "One or more processes are still running"
	cProcessAlreadyRunning      = "Process already running"
	cProcessNotFound            = "Process not found"
)

// Константы указаны в объектах, адрес которых фиксирован всё время работы приложения.
// Ошибку с ошибкой можно сравнивать по телу, по адресу и т.п.
var (
	errSingleton                  = &Error{}
	errUnexpectedError            = err(cUnexpectedError)
	errNotImplemented             = err(cNotImplemented)
	errTypeNotImplemented         = err(cTypeNotImplemented)
	errUnregisterProcessIsRunning = err(cUnregisterProcessIsRunning)
	errRegisteredProcessNotFound  = err(cRegisteredProcessNotFound)
	errDeadlineExceeded           = err(cDeadlineExceeded)
	errProcessesAreStillRunning   = err(cProcessesAreStillRunning)
	errProcessAlreadyRunning      = err(cProcessAlreadyRunning)
	errProcessNotFound            = err(cProcessNotFound)
)

type (
	// Error object of package.
	Error struct{}
	err   string
)

// Error The error built-in interface implementation.
func (e err) Error() string { return string(e) }

// Errors Все ошибки известного состояния, которые могут вернуть функции пакета.
func Errors() *Error { return errSingleton }

// ERRORS:

// UnexpectedError Unexpected error.
func (e *Error) UnexpectedError() error { return &errUnexpectedError }

// NotImplemented Not implemented.
func (e *Error) NotImplemented() error { return &errNotImplemented }

// TypeNotImplemented Type of process is not implemented.
func (e *Error) TypeNotImplemented() error { return &errTypeNotImplemented }

// UnregisterProcessIsRunning The process is current running, you must first stop the process.
func (e *Error) UnregisterProcessIsRunning() error { return &errUnregisterProcessIsRunning }

// RegisteredProcessNotFound Registered process with the specified identifier was not found.
func (e *Error) RegisteredProcessNotFound() error { return &errRegisteredProcessNotFound }

// DeadlineExceeded Deadline exceeded.
func (e *Error) DeadlineExceeded() error { return &errDeadlineExceeded }

// ProcessesAreStillRunning One or more processes are still running.
func (e *Error) ProcessesAreStillRunning() error { return &errProcessesAreStillRunning }

// ProcessAlreadyRunning Process already running.
func (e *Error) ProcessAlreadyRunning() error { return &errProcessAlreadyRunning }

// ProcessNotFound Process not found.
func (e *Error) ProcessNotFound() error { return &errProcessNotFound }
