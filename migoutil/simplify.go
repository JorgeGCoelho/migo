package migoutil

import (
	"github.com/JorgeGCoelho/migo/v3"
	"github.com/JorgeGCoelho/migo/v3/internal/passes/taufunc"
	"github.com/JorgeGCoelho/migo/v3/internal/passes/unused"
)

// SimplifyProgram takes the input Program prog and reduce it
// to a smaller equivalent Program.
//
// It removes functions that reduces to τ, and
// removes call to functions that do not exist.
func SimplifyProgram(prog *migo.Program) *migo.Program {
	if mainmain, hasMM := prog.Function(`"main".main`); hasMM {
		taufunc.Find(prog, taufunc.RemoveExcept(mainmain))
		unused.Remove(prog, mainmain)
	} else {
		taufunc.Find(prog, taufunc.Remove)
	}
	deadcall.Remove(prog)
	return prog
}
