package job // import "gopkg.in/webnice/job.v1/job"

//import "gopkg.in/webnice/debug.v1"
//import "gopkg.in/webnice/log.v2"
import (
	"testing"
)

func TestReset(t *testing.T) {
	var (
		oldContext = singleton.Ctx
		oldPool    = singleton.Pool
		oldWg      = singleton.Wg
	)

	Get().Reset()
	if singleton.CancelFunc == nil {
		t.Fatalf("Reset() error")
	}
	if oldContext == singleton.Ctx {
		t.Fatalf("Reset() error")
	}
	if oldPool == singleton.Pool {
		t.Fatalf("Reset() error")
	}
	if oldWg == singleton.Wg {
		t.Fatalf("Reset() error")
	}
}
