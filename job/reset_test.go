package job

import "testing"

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
