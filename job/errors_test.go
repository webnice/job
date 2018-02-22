package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"testing"
)

func TestErrors(t *testing.T) {
	if ErrorNotImplemented() != errNotImplemented {
		t.Fatalf("The error constants are not correctly defined")
	}
	if ErrorUnregisterProcessIsRunning() != errUnregisterProcessIsRunning {
		t.Fatalf("The error constants are not correctly defined")
	}
	if ErrorRegistredProcessNotFound() != errRegistredProcessNotFound {
		t.Fatalf("The error constants are not correctly defined")
	}
	if ErrorDeadlineExceeded() != errDeadlineExceeded {
		t.Fatalf("The error constants are not correctly defined")
	}
	if ErrorProcessesAreStillRunning() != errProcessesAreStillRunning {
		t.Fatalf("The error constants are not correctly defined")
	}
	if ErrorProcessAlreadyRunning() != errProcessAlreadyRunning {
		t.Fatalf("The error constants are not correctly defined")
	}
	if ErrorProcessNotFound() != errProcessNotFound {
		t.Fatalf("The error constants are not correctly defined")
	}
}
