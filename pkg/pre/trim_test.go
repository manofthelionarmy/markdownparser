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

func TestTrimPreProcessor(t *testing.T) {
	pr, pw := io.Pipe()

	f, err := os.Open("./test.md")
	require.NoError(t, err)
	defer f.Close()

	bw := bytes.NewBuffer(make([]byte, 1024))

	tp := NewTrimProcessor(
		WithSource(f),
		WithTarget(pw),
	)

	wg := sync.WaitGroup{}

	wg.Add(2)
	go func() {
		defer wg.Done()
		// Once we signal to write, we need to immediately read.
		// This can be done via a read via concurrency
		// as long as we read from our pipe reader
		sc := bufio.NewScanner(pr)
		for sc.Scan() {
			fmt.Println(sc.Text())
		}
	}()
	go func() {
		defer wg.Done()
		tp.Process()
	}()
	wg.Wait()
	require.NotZero(t, bw.Len())
	fmt.Println(bw.String())
}
