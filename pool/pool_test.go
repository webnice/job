package pool // import "gopkg.in/webnice/job.v1/pool"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"testing"

	"gopkg.in/webnice/job.v1/types"
)

func TestNew(t *testing.T) {
	var poo *impl
	var ok bool

	poo, ok = New().(*impl)
	if poo == nil {
		t.Fatalf("Error New() function")
	}
	if !ok {
		t.Fatalf("Error New() function")
	}
	if poo.ForkWorkerPool == nil ||
		poo.WorkerPool == nil ||
		poo.TaskPool == nil {
		t.Fatalf("Error init pools")
	}
	if poo.ForkWorkerPool.New == nil ||
		poo.WorkerPool.New == nil ||
		poo.TaskPool.New == nil {
		t.Fatalf("Error init pools")
	}
}

func TestTaskGet(t *testing.T) {
	var poo Interface
	var item1, item2 *types.Task

	poo = New()
	item1 = poo.TaskGet()
	if item1 == nil {
		t.Fatalf("Pool TaskGet error")
	}
	poo.TaskPut(item1)
	item2 = poo.TaskGet()
	if item1 != item2 {
		t.Fatalf("Pool TaskGet error")
	}
	item1 = poo.TaskGet()
	if item1 == item2 {
		t.Fatalf("Pool TaskGet error")
	}
}

func TestWorkerGet(t *testing.T) {
	var poo Interface
	var item1, item2 *types.Worker

	poo = New()
	item1 = poo.WorkerGet()
	if item1 == nil {
		t.Fatalf("Pool WorkerGet error")
	}
	poo.WorkerPut(item1)
	item2 = poo.WorkerGet()
	if item1 != item2 {
		t.Fatalf("Pool WorkerGet error")
	}
	item1 = poo.WorkerGet()
	if item1 == item2 {
		t.Fatalf("Pool WorkerGet error")
	}
}

func TestForkWorkerGet(t *testing.T) {
	var poo Interface
	var item1, item2 *types.ForkWorker

	poo = New()
	item1 = poo.ForkWorkerGet()
	if item1 == nil {
		t.Fatalf("Pool ForkWorkerGet error")
	}
	poo.ForkWorkerPut(item1)
	item2 = poo.ForkWorkerGet()
	if item1 != item2 {
		t.Fatalf("Pool ForkWorkerGet error")
	}
	item1 = poo.ForkWorkerGet()
	if item1 == item2 {
		t.Fatalf("Pool ForkWorkerGet error")
	}
}
