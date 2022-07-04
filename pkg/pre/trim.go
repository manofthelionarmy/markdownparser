package pre

import (
	"bufio"
	"bytes"
	"io"
	"sync"
)

// TrimProcessor trims left and right
type TrimProcessor struct {
	source io.Reader
	target io.Writer
	wg     *sync.WaitGroup
}

// NewTrimProcessor initializes our TrimProcessor
func NewTrimProcessor(opts ...TrimProcessorOpts) *TrimProcessor {
	trimProcessor := &TrimProcessor{}
	for _, f := range opts {
		f(trimProcessor)
	}
	return trimProcessor
}

// It would probably we easier if we make Trim Processor a reader and writer

// Process implements Processor interface to run a preprocessor to clear empty space
func (tp *TrimProcessor) Process() {
	defer tp.wg.Done()

	// PipeWriter needs to send close so that the reader doesn't block!!!
	defer func() {
		if _, ok := tp.target.(*io.PipeWriter); ok {
			tp.target.(*io.PipeWriter).Close()
		}
	}()

	sc := bufio.NewScanner(tp.source)
	for sc.Scan() {
		tp.target.Write(bytes.TrimSpace(sc.Bytes())) // why is this an issue?

		tp.target.Write([]byte("\n"))
	}
}

// SetSource sets the io.reader
func (tp *TrimProcessor) SetSource(r io.Reader) {
	tp.source = r
}

// SetTarget sets the io.Writer
func (tp *TrimProcessor) SetTarget(w io.Writer) {
	tp.target = w
}

// TrimProcessorOpts is options we pass to the initialization of our TrimProcessorOpts
type TrimProcessorOpts func(*TrimProcessor)

// WithSource sets our pipereader to our target
func WithSource(r io.Reader) TrimProcessorOpts {
	return func(tp *TrimProcessor) {
		tp.source = r
	}
}

// WithTarget sets our target
func WithTarget(w io.Writer) TrimProcessorOpts {
	return func(tp *TrimProcessor) {
		tp.target = w
	}
}

// WithWaitGroup sets the WaitGroup
func WithWaitGroup(wg *sync.WaitGroup) TrimProcessorOpts {
	return func(tp *TrimProcessor) {
		tp.wg = wg
	}
}

// WithBufferSize sets the buffer size of our scanner
// func WithBufferSize(size int) TrimProcessorOpts {
// 	return func(tp *TrimProcessor) {
// 		buf := make([]byte, size)
// 		tp.Scanner.Buffer(buf, size)
// 	}
// }
