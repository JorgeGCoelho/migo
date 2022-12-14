package migoutil_test

import (
	"fmt"
	"github.com/JorgeGCoelho/migo/v3"
	"strings"
	"testing"

	"github.com/JorgeGCoelho/migo/v3/migoutil"
	"github.com/JorgeGCoelho/migo/v3/parser"
)

// ErrFuncNotExists is Error if function does not exist in program which is
// should.
type ErrFuncNotExist struct {
	f string
}

func (e *ErrFuncNotExist) Error() string {
	return fmt.Sprintf("Expects %s() in program, but it does not exist", e.f)
}

// ErrFuncExist is Error if function exist in a program which is should not.
type ErrFuncExist struct {
	f string
}

func (e *ErrFuncExist) Error() string {
	return fmt.Sprintf("Expects %s() removed from program, but it exists", e.f)
}

type named struct {
	n string
}

func (n named) Name() string {
	return n.n
}

func (n named) String() string {
	return n.n
}

// Tests SimplifyProgram with simple send/recv/work functions.
//
// def main():
//
//	let ch = newchan ch_instance, 0
//	spawn send(ch)
//	spawn recv(ch)
//	spawn work()
//	recv ch
//	recv ch
//
// def send(sch):
//
//	send sch
//
// def recv(rch):
//
//	recv rch
//	send rch
//
// def work:
//
// main, send, recv should remain after SimplifyProgram
func TestSimplifyProgram(t *testing.T) {
	p := migo.NewProgram()
	mainFunc := migo.NewFunction("main.main")
	mainFunc.AddStmts(
		&migo.NewChanStatement{Name: named{"ch"}, Chan: "ch_instance", Size: 0},
		&migo.SpawnStatement{Name: "send", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}}},
		},
		&migo.SpawnStatement{Name: "recv", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}}},
		},
		&migo.RecvStatement{Chan: "ch"},
		&migo.RecvStatement{Chan: "ch"},
	)
	sendFunc := migo.NewFunction("send")
	sendFunc.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"sch"}})
	sendFunc.AddStmts(
		&migo.SendStatement{Chan: "sch"},
	)
	recvFunc := migo.NewFunction("recv")
	recvFunc.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"rch"}})
	recvFunc.AddStmts(
		&migo.RecvStatement{Chan: "rch"},
		&migo.SendStatement{Chan: "rch"},
	)
	workFunc := migo.NewFunction("work")
	p.AddFunction(mainFunc)
	p.AddFunction(sendFunc)
	p.AddFunction(recvFunc)
	p.AddFunction(workFunc)

	if len(p.Funcs) != 4 {
		t.Errorf("Expects 4 functions in program, but got %d", len(p.Funcs))
	}
	for _, nonempty := range []*migo.Function{mainFunc, sendFunc, recvFunc} {
		if nonempty.IsEmpty() {
			t.Errorf("Expects %s() to be non-empty, but got %t",
				nonempty.Name, nonempty.IsEmpty())
		}
	}
	for _, empty := range []*migo.Function{workFunc} {
		if !empty.IsEmpty() {
			t.Errorf("Expects %s() to be empty, but got %t",
				empty.Name, empty.IsEmpty())
		}
	}
	// These should exist
	for _, exist := range []string{"main.main", "send", "recv", "work"} {
		if _, ok := p.Function(exist); !ok {
			t.Error(&ErrFuncNotExist{f: exist})
		}
	}
	migoutil.SimplifyProgram(p)
	if len(p.Funcs) != 3 {
		t.Errorf("Expects 3 functions in program, but got %d", len(p.Funcs))
	}
	// These should remain
	for _, remain := range []string{"main.main", "send", "recv"} {
		if _, ok := p.Function(remain); !ok {
			t.Error(&ErrFuncNotExist{f: remain})
		}
	}
	// These should be removed
	for _, removed := range []string{"work"} {
		if _, ok := p.Function(removed); ok {
			t.Error(&ErrFuncExist{f: removed})
		}
	}
	s := `
def main.main():
  let ch = newchan ch_instance, 0;
  spawn s(ch);
  spawn r(ch);
  spawn work();
  recv ch;
  recv ch;
def s(sch):
  send sch;
def r(rch):
  recv rch;
  send rch;
def work(): tau;
	`
	prog, err := parser.Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	migoutil.SimplifyProgram(prog)
	if len(prog.Funcs) != 3 {
		t.Errorf("Expects 3 functions in program, but got %d", len(prog.Funcs))
	}
	// These should remain
	for _, remain := range []string{"main.main", "s", "r"} {
		if _, ok := prog.Function(remain); !ok {
			t.Error(&ErrFuncNotExist{f: remain})
		}
	}
	// These should be removed
	for _, removed := range []string{"work"} {
		if _, ok := prog.Function(removed); ok {
			t.Error(&ErrFuncExist{f: removed})
		}
	}
}

// Tests SimplifyProgram with calls to empty functions.
//
// def main():
//
//	let ch = newchan ch_instance, 1
//	call work(ch)
//
// def work(ch):
//
//	call workwork()
//	spawn work$1(ch)
//	call work$2(ch)
//
// def workwork():
//
//	call workworkwork()
//
// def workworkwork():
// def work$1(ch):
// def work$2(ch):
//
//	call work$3(ch)
//
// def work$3(ch):
//
//	recv ch
//	send ch
//
// main, work, work$2, work$3 should remain after SimplifyProgram
func TestSimplifyProgram2(t *testing.T) {
	p := migo.NewProgram()
	mainFunc := migo.NewFunction("main.main")
	mainFunc.AddStmts(
		&migo.NewChanStatement{Name: named{"ch"}, Chan: "ch_instance", Size: 1},
		&migo.CallStatement{Name: "work", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}}},
		},
	)
	workFunc := migo.NewFunction("work")
	workFunc.AddStmts(
		&migo.CallStatement{Name: "workwork"},
		&migo.SpawnStatement{Name: "work$1", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}}},
		},
		&migo.SpawnStatement{Name: "work$2", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}}},
		},
	)
	workworkFunc := migo.NewFunction("workwork")
	workworkFunc.AddStmts(
		&migo.CallStatement{Name: "workworkwork"},
	)
	workworkworkFunc := migo.NewFunction("workworkwork")
	workClosure := migo.NewFunction("work$1")
	workClosure.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}})
	workClosure2 := migo.NewFunction("work$2")
	workClosure2.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}})
	workClosure2.AddStmts(
		&migo.CallStatement{Name: "work$3", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}}},
		},
	)
	workClosure3 := migo.NewFunction("work$3")
	workClosure3.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"ch"}})
	workClosure3.AddStmts(
		&migo.SendStatement{Chan: "ch"},
		&migo.RecvStatement{Chan: "ch"},
	)
	p.AddFunction(mainFunc)
	p.AddFunction(workFunc)
	p.AddFunction(workworkFunc)
	p.AddFunction(workworkworkFunc)
	p.AddFunction(workClosure)
	p.AddFunction(workClosure2)
	p.AddFunction(workClosure3)

	if len(p.Funcs) != 7 {
		t.Errorf("Expects 7 functions in program, but got %d", len(p.Funcs))
	}
	for _, nonEmpty := range []*migo.Function{mainFunc, workFunc, workworkFunc, workClosure2, workClosure3} {
		if nonEmpty.IsEmpty() {
			t.Errorf("Expects %s() to be non-empty, but got %t",
				nonEmpty.Name, nonEmpty.IsEmpty())
		}
	}
	for _, empty := range []*migo.Function{workworkworkFunc, workClosure} {
		if !empty.IsEmpty() {
			t.Errorf("Expects %s() to empty, but got %t",
				empty.Name, empty.IsEmpty())

		}
	}
	for _, exist := range []string{"main.main", "work", "workwork", "workworkwork", "work$1", "work$2", "work$3"} {
		if _, ok := p.Function(exist); !ok {
			t.Error(&ErrFuncNotExist{f: exist})
		}
	}
	migoutil.SimplifyProgram(p)
	if len(p.Funcs) != 4 {
		t.Errorf("Expects 4 functions in program, but got %d", len(p.Funcs))
	}
	// These should remain
	for _, remain := range []string{"main.main", "work", "work$2", "work$3"} {
		if _, ok := p.Function(remain); !ok {
			t.Error(&ErrFuncNotExist{f: remain})
		}
	}
	// These should be removed
	for _, removed := range []string{"workwork", "workworkwork", "work$1"} {
		if _, ok := p.Function(removed); ok {
			t.Error(&ErrFuncExist{f: removed})
		}
	}
	s := `
def main.main():
	let ch = newchan ch_instance, 1;
	call work(ch);
def work(ch):
	call workwork();
	spawn work$1(ch);
	call work$2(ch);
def workwork():
	call workworkwork();
def workworkwork(): tau;
def work$1(ch): tau;
def work$2(ch):
	call work$3(ch);
def work$3(ch):
	recv ch;
	send ch;
`
	prog, err := parser.Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	migoutil.SimplifyProgram(prog)
	if len(prog.Funcs) != 4 {
		t.Errorf("Expects 4 functions in program, but got %d", len(prog.Funcs))
	}
	// These should remain
	for _, remain := range []string{"main.main", "work", "work$2", "work$3"} {
		if _, ok := prog.Function(remain); !ok {
			t.Error(&ErrFuncNotExist{f: remain})
		}
	}
	// These should be removed
	for _, removed := range []string{"workwork", "workworkwork", "work$1"} {
		if _, ok := prog.Function(removed); ok {
			t.Error(&ErrFuncExist{f: removed})
		}
	}
}

// This tests SimplifyProgram from an inconsistent state (main.xx#2 is not defined)
func TestSimplifyProgram3(t *testing.T) {
	s := `def main.main():
    let t0 = newchan main.main.t0_0_0, 0;
    let t1 = newchan main.main.t1_0_0, 0;
    call main.xx(t0, t1);
    spawn main.wait(t0);
    spawn main.wait(t1);
def main.xx(x, y):
    if send x; spawn main.xx(y, x); call main.xx#2(x, y); else call main.xx#2(x, y); endif;
def main.wait(x):
    call main.wait#1(x);
def main.wait#1(x):
    recv x;
    call main.wait#1(x);`
	expect := `def main.main():
    let t0 = newchan main.main.t0_0_0, 0;
    let t1 = newchan main.main.t1_0_0, 0;
    call main.xx(t0, t1);
    spawn main.wait(t0);
    spawn main.wait(t1);
def main.xx(x, y):
    if send x; spawn main.xx(y, x); else tau; endif;
def main.wait(x):
    call main.wait#1(x);
def main.wait#1(x):
    recv x;
    call main.wait#1(x);`
	r := strings.NewReader(s)
	parsed, err := parser.Parse(r)
	migoutil.SimplifyProgram(parsed)
	if err != nil {
		t.Error(err)
	}
	if strings.TrimSpace(parsed.String()) != expect {
		t.Errorf("Expects main.xx#2 calls to be removed\n--- expect ---\n%s\n--- got ---\n%s\n",
			expect, parsed.String())
	}
	prog, err := parser.Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	migoutil.SimplifyProgram(prog)
	if strings.TrimSpace(prog.String()) != expect {
		t.Errorf("Expects main.xx#2 calls to be removed\n--- expect ---\n%s\n--- got ---\n%s\n",
			expect, prog.String())
	}
}

// Tests SimplifyProgram in the presence of mem ops.
func TestSimplfyProgramMemOp(t *testing.T) {
	p := migo.NewProgram()
	mainFunc := migo.NewFunction("main.main")
	mainFunc.AddStmts(
		&migo.NewChanStatement{Name: named{"ch"}, Chan: "ch_instance", Size: 0},
		&migo.NewMem{Name: named{"mem0"}},
		&migo.SpawnStatement{Name: "send", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}}},
		},
		&migo.SpawnStatement{Name: "recv", Params: []*migo.Parameter{
			&migo.Parameter{Caller: named{"ch"}}},
		},
		&migo.RecvStatement{Chan: "ch"},
		&migo.RecvStatement{Chan: "ch"},
		&migo.MemWrite{Name: "mem0"},
	)
	sendFunc := migo.NewFunction("send")
	sendFunc.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"sch"}})
	sendFunc.AddStmts(
		&migo.SendStatement{Chan: "sch"},
		&migo.MemWrite{Name: "mem0"},
	)
	recvFunc := migo.NewFunction("recv")
	recvFunc.AddParams(&migo.Parameter{Caller: named{"ch"}, Callee: named{"rch"}})
	recvFunc.AddStmts(
		&migo.RecvStatement{Chan: "rch"},
		&migo.SendStatement{Chan: "rch"},
	)
	workFunc := migo.NewFunction("work")
	workFunc.AddStmts(
		&migo.MemRead{Name: "mem0"},
		&migo.MemWrite{Name: "mem0"},
	)
	emptyFunc := migo.NewFunction("empty")
	p.AddFunction(mainFunc)
	p.AddFunction(sendFunc)
	p.AddFunction(recvFunc)
	p.AddFunction(workFunc)
	p.AddFunction(emptyFunc)

	if len(p.Funcs) != 5 {
		t.Errorf("Expects 5 functions in program, but got %d", len(p.Funcs))
	}
	for _, nonempty := range []*migo.Function{mainFunc, sendFunc, recvFunc, workFunc} {
		if nonempty.IsEmpty() {
			t.Errorf("Expects %s() to be non-empty, but got %t",
				nonempty.Name, nonempty.IsEmpty())
		}
	}
	for _, empty := range []*migo.Function{emptyFunc} {
		if !empty.IsEmpty() {
			t.Errorf("Expects %s() to be empty, but got %t",
				empty.Name, empty.IsEmpty())
		}
	}
	// These should exist
	for _, exist := range []string{"main.main", "send", "recv", "work"} {
		if _, ok := p.Function(exist); !ok {
			t.Error(&ErrFuncNotExist{f: exist})
		}
	}
	migoutil.SimplifyProgram(p)
	if len(p.Funcs) != 4 {
		t.Errorf("Expects 4 functions in program, but got %d", len(p.Funcs))
	}
	// These should remain
	for _, remain := range []string{"main.main", "send", "recv", "work"} {
		if _, ok := p.Function(remain); !ok {
			t.Error(&ErrFuncNotExist{f: remain})
		}
	}
	// These should be removed
	for _, removed := range []string{"empty"} {
		if _, ok := p.Function(removed); ok {
			t.Error(&ErrFuncExist{f: removed})
		}
	}
	s := `
def main.main():
  let ch = newchan ch_instance, 0;
  spawn s(ch);
  spawn r(ch);
  spawn work();
  recv ch;
  recv ch;
def s(sch):
  send sch;
def r(rch):
  recv rch;
  send rch;
def work(): tau;
	`
	prog, err := parser.Parse(strings.NewReader(s))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	migoutil.SimplifyProgram(prog)
	if len(prog.Funcs) != 3 {
		t.Errorf("Expects 3 functions in program, but got %d", len(prog.Funcs))
	}
	// These should remain
	for _, remain := range []string{"main.main", "s", "r"} {
		if _, ok := prog.Function(remain); !ok {
			t.Error(&ErrFuncNotExist{f: remain})
		}
	}
	// These should be removed
	for _, removed := range []string{"work"} {
		if _, ok := prog.Function(removed); ok {
			t.Error(&ErrFuncExist{f: removed})
		}
	}
}
