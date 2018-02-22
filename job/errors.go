package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import "fmt"

var (
	errNotImplemented             = fmt.Errorf("Not implemented")
	errUnregisterProcessIsRunning = fmt.Errorf("The process is current running, you must first stop the process")
	errRegistredProcessNotFound   = fmt.Errorf("Registred process with the specified identifier was not found")
	errDeadlineExceeded           = fmt.Errorf("Deadline exceeded")
	errProcessesAreStillRunning   = fmt.Errorf("One or more processes are still running")
	errProcessAlreadyRunning      = fmt.Errorf("Process already running")
	errProcessNotFound            = fmt.Errorf("Process not found")
)

// ErrorNotImplemented Not implemented
func ErrorNotImplemented() error { return errNotImplemented }

// ErrorUnregisterProcessIsRunning The process is current running, you must first stop the process
func ErrorUnregisterProcessIsRunning() error { return errUnregisterProcessIsRunning }

// ErrorRegistredProcessNotFound Registred process with the specified identifier was not found
func ErrorRegistredProcessNotFound() error { return errRegistredProcessNotFound }

// ErrorDeadlineExceeded Deadline exceeded
func ErrorDeadlineExceeded() error { return errDeadlineExceeded }

// ErrorProcessesAreStillRunning One or more processes are still running
func ErrorProcessesAreStillRunning() error { return errProcessesAreStillRunning }

// ErrorProcessAlreadyRunning Process already running
func ErrorProcessAlreadyRunning() error { return errProcessAlreadyRunning }

// ErrorProcessNotFound Process not found
func ErrorProcessNotFound() error { return errProcessNotFound }
