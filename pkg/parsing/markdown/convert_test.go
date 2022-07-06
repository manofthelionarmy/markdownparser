package markdown

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	for scenario, f := range map[string]func(*testing.T){
		"testHeaders":       testHeaders,
		"testParagraph":     testParagraph,
		"testLinks":         testLinks,
		"testNestedLinks":   testNestedLinks,
		"testHasNestedLink": testHasNestedLink,
		"testUnorderedList": testUnorderedList,
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

func testNestedLinks(t *testing.T) {
	require.Equal(t,
		"<p>This has a nested <a href=\"https://url.com\" alt=\"Optional Alt\">Link Text</a></p>\n",
		string(convert("This has a nested [Link Text](https://url.com \"Optional Alt\")")),
	)
	require.Equal(t,
		"<h1>This has a nested <a href=\"https://url.com\" alt=\"Optional Alt\">Link Text</a></h1>\n",
		string(convert("# This has a nested [Link Text](https://url.com \"Optional Alt\")")),
	)
	markdownList := "* [Link Text](https://url.com \"Optional Alt\")\n"
	markdownList = markdownList + markdownList + markdownList

	listHTML := "<li><a href=\"https://url.com\" alt=\"Optional Alt\">Link Text</a></li>\n"
	listHTML = "<ul>\n" + listHTML + listHTML + listHTML + "</ul>\n"
	require.Equal(t,
		listHTML,
		string(convert(markdownList)),
	)
}

func testHasNestedLink(t *testing.T) {
	require.True(t, hasNestedLink("This has a nested [Link Text](https://url.com \"Optional Alt\")"))
}

func testUnorderedList(t *testing.T) {
	markdownList := "* Fruits\n" + "* Milk\n" + "* Bread\n"

	htmlList := "<ul>\n" +
		"<li>Fruits</li>\n" +
		"<li>Milk</li>\n" +
		"<li>Bread</li>\n" +
		"</ul>\n"

	require.Equal(
		t,
		"<ul>\n"+"<li></li>\n"+"</ul>\n",
		string(convert("* \n")),
	)

	require.Equal(
		t,
		htmlList,
		string(convert(markdownList)),
	)
}
