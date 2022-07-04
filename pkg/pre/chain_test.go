package pre

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainProcessor(t *testing.T) {
	for scenario, f := range map[string]func(*testing.T){
		"testChainProcessors": testChainProcessors,
	} {
		t.Run(scenario, f)
	}
}

func TestChain(t *testing.T) {
	f, err := os.Open("test.md")
	require.NoError(t, err)

	trimProcessor := &TrimProcessor{}
	upperCase := &upperCaseProcessor{}
	lowerCase := &lowerCaseProcessor{}

	ChainProcessors(f, []Processor{trimProcessor, upperCase, lowerCase}...)
	require.Equal(t, f, trimProcessor.source)
	require.NotEqual(t, trimProcessor.target, upperCase.source)
	require.NotEqual(t, upperCase.target, lowerCase.source)
	require.NotEqual(t, os.Stdout, lowerCase.target)
}

func testChainProcessors(t *testing.T) {
	f, err := os.Open("test.md")
	require.NoError(t, err)

	wg := &sync.WaitGroup{}
	trimProcessor := NewTrimProcessor(
		WithWaitGroup(wg),
	)
	upperCase := &upperCaseProcessor{wg: wg}
	lowerCase := &lowerCaseProcessor{wg: wg}

	r := ChainProcessors(f, []Processor{trimProcessor, upperCase}...)

	// After chain, we need one more pipe
	lr, lw := io.Pipe()
	lowerCase.source = r
	lowerCase.target = lw

	wg.Add(4)

	go trimProcessor.Process()

	go upperCase.Process()

	go lowerCase.Process()

	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, lr)
	}()

	wg.Wait()
}

type upperCaseProcessor struct {
	*io.PipeWriter
	*io.PipeReader
	target io.Writer
	source io.Reader
	wg     *sync.WaitGroup
}

func (u *upperCaseProcessor) Process() {
	defer u.wg.Done()
	defer func() {
		if _, ok := u.target.(*io.PipeWriter); ok {
			u.target.(*io.PipeWriter).Close()
		}
	}()

	sc := bufio.NewScanner(u.source)
	for sc.Scan() {
		u.target.Write(bytes.ToUpper(sc.Bytes()))
		u.target.Write([]byte("\n"))
	}
}

func (u *upperCaseProcessor) SetSource(r io.Reader) {
	u.source = r
}

func (u *upperCaseProcessor) SetTarget(w io.Writer) {
	u.target = w
}

type lowerCaseProcessor struct {
	*io.PipeWriter
	*io.PipeReader
	target io.Writer
	source io.Reader
	wg     *sync.WaitGroup
}

func newLowerCaseProcessor() *lowerCaseProcessor {
	return &lowerCaseProcessor{}
}

func (l *lowerCaseProcessor) Process() {
	defer l.wg.Done()
	defer func() {
		if _, ok := l.target.(*io.PipeWriter); ok {
			l.target.(*io.PipeWriter).Close()
		}
	}()
	sc := bufio.NewScanner(l.source)
	for sc.Scan() {
		l.target.Write(bytes.ToLower(sc.Bytes()))
		l.target.Write([]byte("\n"))
	}
}

func (l *lowerCaseProcessor) SetSource(r io.Reader) {
	l.source = r
}

func (l *lowerCaseProcessor) SetTarget(w io.Writer) {
	l.target = w
}
