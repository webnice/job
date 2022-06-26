package job

import (
	"strings"
	"testing"
)

func TestErrors(t *testing.T) {
	const err = `The error constants are not correctly defined`

	if Errors().UnexpectedError() != &errUnexpectedError {
		t.Fatalf(err)
	}
	if Errors().NotImplemented() != &errNotImplemented {
		t.Fatalf(err)
	}
	if Errors().TypeNotImplemented() != &errTypeNotImplemented {
		t.Fatalf(err)
	}
	if Errors().UnregisterProcessIsRunning() != &errUnregisterProcessIsRunning {
		t.Fatalf(err)
	}
	if Errors().RegisteredProcessNotFound() != &errRegisteredProcessNotFound {
		t.Fatalf(err)
	}
	if Errors().DeadlineExceeded() != &errDeadlineExceeded {
		t.Fatalf(err)
	}
	if Errors().ProcessesAreStillRunning() != &errProcessesAreStillRunning {
		t.Fatalf(err)
	}
	if Errors().ProcessAlreadyRunning() != &errProcessAlreadyRunning {
		t.Fatalf(err)
	}
	if Errors().ProcessNotFound() != &errProcessNotFound {
		t.Fatalf(err)
	}
	if strings.Compare(Errors().UnexpectedError().Error(), cUnexpectedError) != 0 {
		t.Fatalf(`Errors implemented not correctly`)
	}
}
