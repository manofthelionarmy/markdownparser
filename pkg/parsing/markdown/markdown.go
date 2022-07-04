package markdown

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/manofthelionarmy/markdownparser/pkg/pre"
)

// Config is the configuration for our markdown parser
type Config struct {
	Concurrency int
}

// Parser is our markdown parser
type Parser struct {
	source        io.Reader
	target        io.Writer
	mainTarget    io.Writer
	preprocessors []pre.Processor
	wg            *sync.WaitGroup
}

// NewMarkdownParser returns a new markdown parser
func NewMarkdownParser(cfg *Config, opts ...ParsingOpts) *Parser {
	markdownParser := &Parser{}
	for _, f := range opts {
		f(markdownParser)
	}
	return markdownParser
}

// Parse begins parsing the markdown
func (md *Parser) Parse() {

	// Chain all of our preprocessors and return the pipe reader for the last io.Pipe() created
	r := pre.ChainProcessors(md.source, md.preprocessors...)

	// need to set up the pipe for our parser
	markDownPr, markDownWr := io.Pipe()
	md.source = r // the source is the reader the last preprocessor writes to
	md.mainTarget = md.target
	md.target = markDownWr

	wg := sync.WaitGroup{}

	for _, p := range md.preprocessors {
		wg.Add(1)
		go func(p pre.Processor) {
			defer wg.Done()
			p.Process()
		}(p)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() {
			if _, ok := md.target.(*io.PipeWriter); ok {
				md.target.(*io.PipeWriter).Close()
			}
		}()
		sc := bufio.NewScanner(md.source)
		for sc.Scan() {
			if h1Regex.Match(sc.Bytes()) {
				text := strings.Split(sc.Text(), "#")[1]

				text = strings.TrimSpace(text)
				text = "<h1>" + text + "</h1>"

				md.target.Write([]byte(text))
				md.target.Write([]byte("\n"))
			} else if h2Regex.Match(sc.Bytes()) {
				text := strings.Split(sc.Text(), "##")[1]

				text = strings.TrimSpace(text)
				text = "<h2>" + text + "</h2>"

				md.target.Write([]byte(text))
				md.target.Write([]byte("\n"))
			}
		}
	}()

	go func() {
		defer wg.Done()
		if _, err := io.Copy(md.mainTarget, markDownPr); err != nil {
			log.Fatal(err)
		}
	}()
	wg.Wait()
}

// ParsingOpts is a function that sets specified behavior for our parser
type ParsingOpts func(*Parser)

// WithPipeWriter pipes our parser generated output to a target stream
// func WithPipeWriter(w *io.PipeWriter) ParsingOpts {
// 	return func(mp *Parser) {
// 		mp.PipeWriter = w
// 	}
// }

// WithSource sets the read half of a pipe
func WithSource(r io.Reader) ParsingOpts {
	// io.Reader differs from PipeReader
	// we can do an io.copy of our scanner to the pipe reader?
	return func(mp *Parser) {
		mp.source = r
	}
}

// WithTarget sets the destination
func WithTarget(w io.Writer) ParsingOpts {
	return func(p *Parser) {
		p.target = w
	}
}

// WithPreprocessor does some stuff to modify the data prior to parsing
func WithPreprocessor(processor pre.Processor) ParsingOpts {
	return func(p *Parser) {
		p.preprocessors = append(p.preprocessors, processor)
	}
}

// WithWaitGroup sets the WaitGroup
func WithWaitGroup(wg *sync.WaitGroup) ParsingOpts {
	return func(p *Parser) {
		p.wg = wg
	}
}
