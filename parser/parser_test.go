package parser

import (
	"github.com/JorgeGCoelho/migo/v3"
	"strings"
	"testing"
)

func TestParser(t *testing.T) {
	want := `def main():
    let ch = newchan T, 0;
    send ch;
`
	p, err := Parse(strings.NewReader(`   def main(): let ch = newchan T, 0;
	send ch;   `))
	if err != nil {
		t.Errorf("cannot parse: %v", err)
	}
	if got := p.String(); want != got {
		t.Errorf("unexpected parsed migo, want:\n%sgot:\n%s", want, got)
	}
}

func TestParseMem(t *testing.T) {
	s := `def main(): letmem x; read x; spawn f(); write x;
    def f(): read x; write x;`
	parsed, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if want, got := 2, len(parsed.Funcs); want != got {
		t.Errorf("expected %d functions but got %d", want, got)
	}
	fn, found := parsed.Function("main")
	if !found {
		t.Error("cannot find main function")
	}
	if want, got := 4, len(fn.Stmts); want != got {
		t.Errorf("expected %d statements but got %d", want, got)
	}
	stmt0, ok := fn.Stmts[0].(*migo.NewMem)
	if !ok {
		t.Errorf("expecting letmem statement but got %v", fn.Stmts[0])
		t.FailNow()
	}
	if stmt0.Name.Name() != "x" {
		t.Errorf("expected letmem x but got letmem %s", stmt0.Name)
	}
	stmt1, ok := fn.Stmts[1].(*migo.MemRead)
	if !ok {
		t.Errorf("expecting read statement but got %v", fn.Stmts[1])
		t.FailNow()
	}
	if stmt1.Name != "x" {
		t.Errorf("expected read x but got letmem %s", stmt1.Name)
	}
	stmt3, ok := fn.Stmts[3].(*migo.MemWrite)
	if !ok {
		t.Errorf("expecting write statement but got %v", fn.Stmts[3])
		t.FailNow()
	}
	if stmt3.Name != "x" {
		t.Errorf("expected write x but got letmem %s", stmt3.Name)
	}
}

func TestParseMutex(t *testing.T) {
	s := `def main(): letsync a mutex; lock a; unlock a;`
	parsed, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fn, found := parsed.Function("main")
	if !found {
		t.Error("cannot find main function")
	}
	if want, got := 3, len(fn.Stmts); want != got {
		t.Errorf("expected %d statements but got %d", want, got)
	}
	stmt0, ok := fn.Stmts[0].(*migo.NewSyncMutex)
	if !ok {
		t.Errorf("expecting letsync statement but got %v", fn.Stmts[0])
		t.FailNow()
	}
	if stmt0.Name.Name() != "a" {
		t.Errorf("expected letsync a mutex but got %v", fn.Stmts[0])
	}
	stmt1, ok := fn.Stmts[1].(*migo.SyncMutexLock)
	if !ok {
		t.Errorf("expecting lock statement but got %v", fn.Stmts[1])
		t.FailNow()
	}
	if stmt1.Name != "a" {
		t.Errorf("expected lock a but got %v", fn.Stmts[1])
	}
	stmt2, ok := fn.Stmts[2].(*migo.SyncMutexUnlock)
	if !ok {
		t.Errorf("expecting unlock statement but got %v", fn.Stmts[2])
		t.FailNow()
	}
	if stmt2.Name != "a" {
		t.Errorf("expected unlock a but got %v", fn.Stmts[2])
	}
}

func TestParseRWMutex(t *testing.T) {
	s := `def main(): letsync a rwmutex; lock a; unlock a; rlock a; runlock a;`
	parsed, err := Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fn, found := parsed.Function("main")
	if !found {
		t.Error("cannot find main function")
	}
	if want, got := 5, len(fn.Stmts); want != got {
		t.Errorf("expected %d statements but got %d", want, got)
	}
	stmt0, ok := fn.Stmts[0].(*migo.NewSyncRWMutex)
	if !ok {
		t.Errorf("expecting letsync statement but got %v", fn.Stmts[0])
		t.FailNow()
	}
	if stmt0.Name.Name() != "a" {
		t.Errorf("expected letsync a rwmutex but got %v", fn.Stmts[0])
	}
	stmt1, ok := fn.Stmts[1].(*migo.SyncMutexLock)
	if !ok {
		t.Errorf("expecting lock statement but got %v", fn.Stmts[1])
		t.FailNow()
	}
	if stmt1.Name != "a" {
		t.Errorf("expected lock a but got %v", fn.Stmts[1])
	}
	stmt2, ok := fn.Stmts[2].(*migo.SyncMutexUnlock)
	if !ok {
		t.Errorf("expecting unlock statement but got %v", fn.Stmts[2])
		t.FailNow()
	}
	if stmt2.Name != "a" {
		t.Errorf("expected unlock a but got %v", fn.Stmts[2])
	}
	stmt3, ok := fn.Stmts[3].(*migo.SyncRWMutexRLock)
	if !ok {
		t.Errorf("expecting rlock statement but got %v", fn.Stmts[3])
		t.FailNow()
	}
	if stmt3.Name != "a" {
		t.Errorf("expected rlock a but got %v", fn.Stmts[3])
	}
	stmt4, ok := fn.Stmts[4].(*migo.SyncRWMutexRUnlock)
	if !ok {
		t.Errorf("expecting runlock statement but got %v", fn.Stmts[4])
		t.FailNow()
	}
	if stmt4.Name != "a" {
		t.Errorf("expected runlock a but got %v", fn.Stmts[4])
	}
}
