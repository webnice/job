package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"testing"
)

func TestType(t *testing.T) {
	if Type(`task`).String() != `task` {
		t.Errorf("Error type task")
		return
	}
	if Type(`worker`).String() != `worker` {
		t.Errorf("Error type worker")
		return
	}
	if Type(`forkworker`).String() != `forkworker` {
		t.Errorf("Error type fork worker")
		return
	}
}
