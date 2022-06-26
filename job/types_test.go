package job

import (
	"testing"

	jobTypes "github.com/webnice/job/types"
)

func TestType(t *testing.T) {
	if jobTypes.Type(`task`).String() != `task` {
		t.Errorf("Error type task")
		return
	}
	if jobTypes.Type(`worker`).String() != `worker` {
		t.Errorf("Error type worker")
		return
	}
	if jobTypes.Type(`forkworker`).String() != `forkworker` {
		t.Errorf("Error type fork worker")
		return
	}
}
