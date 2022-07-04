package markdown

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	for scenario, f := range map[string]func(*testing.T){
		"testHeaders":   testHeaders,
		"testParagraph": testParagraph,
		"testLinks":     testLinks,
	} {
		t.Run(scenario, f)
	}
}

func testHeaders(t *testing.T) {
	require.Equal(t, []byte("<h1>hello</h1>\n"), convert("# hello"))
	require.Equal(t, []byte("<h2>hello</h2>\n"), convert("## hello"))
	require.Equal(t, []byte("<h3>hello</h3>\n"), convert("### hello"))
	require.Equal(t, []byte("<h4>hello</h4>\n"), convert("#### hello"))
	require.Equal(t, []byte("<h5>hello</h5>\n"), convert("##### hello"))
	require.Equal(t, []byte("<h6>hello</h6>\n"), convert("###### hello"))
}

func testParagraph(t *testing.T) {
	// TODO: a nice to have, post processor to do a word wrap on paragraphs
	require.Equal(t,
		"<p>This is my whole paragraph ###Heading should be skipped.</p>\n",
		string(convert("This is my whole paragraph ###Heading should be skipped.")),
	)
}

func testLinks(t *testing.T) {
	t.Errorf("Implement test")
}
