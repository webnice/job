package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"testing"

	jobTypes "gopkg.in/webnice/job.v1/types"
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
