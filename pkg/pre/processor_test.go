package pre

import (
	"bufio"
	"bytes"
	"fmt"
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
	// lowerCase := &lowerCaseProcessor{}

	ChainProcessors(f, []Processor{trimProcessor, upperCase}...)
	require.Equal(t, f, trimProcessor.source)
	require.NotEqual(t, trimProcessor.target, upperCase.source)
	// require.NotEqual(t, upperCase.target, lowerCase.source)
	// require.Equal(t, os.Stdout, lowerCase.target)
	// require.Equal(t, os.Stdout, upperCase.target)
}

func testChainProcessors(t *testing.T) {
	// pr, pw := io.Pipe()

	f, err := os.Open("test.md")
	require.NoError(t, err)

	trimProcessor := &TrimProcessor{}
	upperCase := &upperCaseProcessor{}
	lowerCase := &lowerCaseProcessor{}

	// lwPr, lwPw := io.Pipe()

	// trimProcessor.source = f
	// trimProcessor.target = pw

	// upperCase.source = pr
	// upperCase.target = lwPw

	// lowerCase.source = lwPr
	// lowerCase.target = os.Stdout
	r := ChainProcessors(f, []Processor{trimProcessor, upperCase}...)

	// After chain, we need one more pipe
	lr, lw := io.Pipe()
	lowerCase.source = r
	lowerCase.target = lw

	wg := sync.WaitGroup{}
	wg.Add(4)

	go func() {
		defer wg.Done()
		defer fmt.Println("chaining")
		trimProcessor.Process()
	}()

	go func() {
		defer wg.Done()
		upperCase.Process()
	}()

	go func() {
		defer wg.Done()
		lowerCase.Process()
	}()

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
}

func (u *upperCaseProcessor) Process() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
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
	}()
	wg.Wait()
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
}

func newLowerCaseProcessor() *lowerCaseProcessor {
	return &lowerCaseProcessor{}
}

func (l *lowerCaseProcessor) Process() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
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
	}()
	defer wg.Wait()
}

func (l *lowerCaseProcessor) SetSource(r io.Reader) {
	l.source = r
}

func (l *lowerCaseProcessor) SetTarget(w io.Writer) {
	l.target = w
}
