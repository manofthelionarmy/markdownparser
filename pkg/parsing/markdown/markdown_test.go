package markdown

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"testing"

	"github.com/manofthelionarmy/markdownparser/pkg/pre"
	"github.com/stretchr/testify/require"
)

func TestParsing(t *testing.T) {
	for scenario, f := range map[string]func(t *testing.T){
		"testNewParser":               testNewParser,
		"testPipeTrimToMarkDownParse": testPipeTrimToMarkDownParse,
	} {
		t.Run(scenario, f)
	}
}

func testNewParser(t *testing.T) {
	cfg := &Config{}
	parser := NewMarkdownParser(cfg)
	require.NotNil(t, parser)
}

func testPipeTrimToMarkDownParse(t *testing.T) {
	f, err := os.Open("test.md")
	require.NoError(t, err)
	defer f.Close()

	trimProcessor := pre.NewTrimProcessor()

	cfg := &Config{}

	buf := make([]byte, 1024)
	bw := bytes.NewBuffer(buf)

	markDownParser := NewMarkdownParser(
		cfg,
		WithSource(f),
		WithPreprocessor(trimProcessor),
		WithTarget(bw),
	)

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		markDownParser.Parse()
	}()

	wg.Wait()
	require.NotZero(t, bw.Len())
	fmt.Println(bw.String())
}
