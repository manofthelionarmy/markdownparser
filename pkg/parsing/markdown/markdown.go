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
	writer        io.Writer
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
	md.mainTarget = md.writer
	md.writer = markDownWr

	for _, p := range md.preprocessors {
		md.wg.Add(1)
		go p.Process()
	}

	md.wg.Add(1)
	go func() {
		defer md.wg.Done()
		defer func() {
			if _, ok := md.writer.(*io.PipeWriter); ok {
				md.writer.(*io.PipeWriter).Close()
			}
		}()
		sc := bufio.NewScanner(md.source)
		for sc.Scan() {
			// TODO: make this testable, break this out into a function
			if h1Regex.Match(sc.Bytes()) {
				text := strings.Split(sc.Text(), "#")[1]

				text = strings.TrimSpace(text)
				text = "<h1>" + text + "</h1>"

				md.writer.Write([]byte(text + "\n"))
			} else if h2Regex.Match(sc.Bytes()) {
				text := strings.Split(sc.Text(), "##")[1]

				text = strings.TrimSpace(text)
				text = "<h2>" + text + "</h2>"

				md.writer.Write([]byte(text + "\n"))
			}
		}
	}()

	md.wg.Add(1)
	go func() {
		defer md.wg.Done()
		if _, err := io.Copy(md.mainTarget, markDownPr); err != nil {
			log.Fatal(err)
		}
	}()
	md.wg.Wait()
}

// ParsingOpts is a function that sets specified behavior for our parser
type ParsingOpts func(*Parser)

// WithSource sets the read half of a pipe
func WithSource(r io.Reader) ParsingOpts {
	return func(mp *Parser) {
		mp.source = r
	}
}

// WithTarget sets the destination
func WithTarget(w io.Writer) ParsingOpts {
	return func(p *Parser) {
		p.writer = w
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
