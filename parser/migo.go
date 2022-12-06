package parser

import "github.com/JorgeGCoelho/migo/v3"

// Helper functions for yacc parser
// These functions wrap MiGo AST

func sendStmt(ch string) *migo.SendStatement {
	return &migo.SendStatement{Chan: ch}
}

func recvStmt(ch string) *migo.RecvStatement {
	return &migo.RecvStatement{Chan: ch}
}

func tauStmt() *migo.TauStatement {
	return &migo.TauStatement{}
}

func newchanStmt(name, ch string, size int) migo.Statement {
	return &migo.NewChanStatement{
		Name: &plainNamedVar{s: name},
		Chan: ch,
		Size: int64(size),
	}
}

func newMutex(name string) *migo.NewSyncMutex {
	return &migo.NewSyncMutex{Name: name}
}

func lockStmt(name string) *migo.SyncMutexLock {
	return &migo.SyncMutexLock{Name: name}
}

func unlockStmt(name string) *migo.SyncMutexUnlock {
	return &migo.SyncMutexUnlock{Name: name}
}

func newRWMutex(name string) *migo.NewSyncRWMutex {
	return &migo.NewSyncRWMutex{Name: name}
}

func rlockStmt(name string) *migo.SyncRWMutexRLock {
	return &migo.SyncRWMutexRLock{Name: name}
}

func runlockStmt(name string) *migo.SyncRWMutexRUnlock {
	return &migo.SyncRWMutexRUnlock{Name: name}
}

func readStmt(name string) *migo.MemRead {
	return &migo.MemRead{Name: name}
}

func writeStmt(name string) *migo.MemWrite {
	return &migo.MemWrite{Name: name}
}

func newmemStmt(name string) *migo.NewMem {
	return &migo.NewMem{Name: name}
}

func closeStmt(ch string) *migo.CloseStatement {
	return &migo.CloseStatement{Chan: ch}
}

func callStmt(fn string, params []*migo.Parameter) *migo.CallStatement {
	return &migo.CallStatement{Name: fn, Params: params}
}

func spawnStmt(fn string, params []*migo.Parameter) *migo.SpawnStatement {
	return &migo.SpawnStatement{Name: fn, Params: params}
}

func params(p ...*migo.Parameter) []*migo.Parameter {
	return p
}

func plainParam(name string) *migo.Parameter {
	return &migo.Parameter{Caller: &plainNamedVar{s: name}, Callee: &plainNamedVar{s: name}}
}

func ifStmt(iftrue, iffalse []migo.Statement) *migo.IfStatement {
	return &migo.IfStatement{Then: iftrue, Else: iffalse}
}

func selectStmt(cases [][]migo.Statement) *migo.SelectStatement {
	return &migo.SelectStatement{Cases: cases}
}

func cases(c ...[]migo.Statement) [][]migo.Statement {
	return c
}

func stmts(s ...migo.Statement) []migo.Statement {
	return s
}
