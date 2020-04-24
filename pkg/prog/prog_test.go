package prog

import (
	"os"
	"testing"

	. "github.com/elves/elvish/pkg/prog/progtest"
)

func TestBadFlag(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	exit := Run(f.Fds(), Elvish("-bad-flag"))

	TestError(t, f, exit, "flag provided but not defined: -bad-flag")
}

func TestHelp(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	Run(f.Fds(), Elvish("-help"))

	f.TestOutSnippet(t, 1, "Usage: elvish [flags] [script]")
}

func TestNoProgram(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	exit := Run(f.Fds(), Elvish(), testProgram{}, testProgram{})

	TestError(t, f, exit, "program bug: no suitable subprogram")
}

func TestGoodProgram(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	Run(f.Fds(), Elvish(), testProgram{},
		testProgram{shouldRun: true, writeOut: "program 2"})

	f.TestOut(t, 1, "program 2")
}

func TestPreferEarlierProgram(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	Run(f.Fds(), Elvish(),
		testProgram{shouldRun: true, writeOut: "program 1"},
		testProgram{shouldRun: true, writeOut: "program 2"})

	f.TestOut(t, 1, "program 1")
}

func TestBadUsageError(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	exit := Run(f.Fds(), Elvish(),
		testProgram{shouldRun: true, returnErr: BadUsage("lorem ipsum")})

	TestError(t, f, exit, "lorem ipsum")
	f.TestOutSnippet(t, 2, "Usage:")
}

func TestExitError(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	exit := Run(f.Fds(), Elvish(),
		testProgram{shouldRun: true, returnErr: Exit(3)})

	if exit != 3 {
		t.Errorf("exit = %v, want 3", exit)
	}
	f.TestOut(t, 2, "")
}

func TestExitError_0(t *testing.T) {
	f := Setup()
	defer f.Cleanup()

	exit := Run(f.Fds(), Elvish(),
		testProgram{shouldRun: true, returnErr: Exit(0)})

	if exit != 0 {
		t.Errorf("exit = %v, want 0", exit)
	}
}

type testProgram struct {
	shouldRun bool
	writeOut  string
	returnErr error
}

func (p testProgram) ShouldRun(*Flags) bool { return p.shouldRun }

func (p testProgram) Run(fds [3]*os.File, _ *Flags, args []string) error {
	fds[1].WriteString(p.writeOut)
	return p.returnErr
}