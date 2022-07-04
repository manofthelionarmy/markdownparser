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
		"testImages":    testImages,
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
	require.Equal(t,
		"<a href=\"https://url.com\" alt=\"Optional Alt\">Link Text</a>\n",
		string(convert("[Link Text](https://url.com \"Optional Alt\")")),
	)

	// A lot of space
	require.Equal(t,
		"<a href=\"https://url.com\" alt=\"Optional Alt\">Link Text</a>\n",
		string(convert(`[Link Text](https://url.com          "Optional Alt")`)),
	)
	// No optional alt
	require.Equal(t,
		"<a href=\"https://url.com\">Link Text</a>\n",
		string(convert("[Link Text](https://url.com)")),
	)
}

func testImages(t *testing.T) {
	require.Equal(t,
		`<img src="path/to/img.png" alt="Alternative text"/>`+"\n",
		string(convert(`![Alternative text](path/to/img.png "Text")`)),
	)
	require.Equal(t,
		`<img src="/assets/img/MarineGEO_logo.png" alt="MarineGEO circle logo"/>`+"\n",
		string(convert(`![MarineGEO circle logo](/assets/img/MarineGEO_logo.png "MarineGEO logo")`)),
	)
}
