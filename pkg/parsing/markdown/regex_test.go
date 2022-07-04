package markdown

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarkdownRegex(t *testing.T) {
	for scenario, f := range map[string]func(*testing.T){
		"testHeader1": testHeader1,
		"testHeader2": testHeader2,
		// "testHeader3": testHeader3,
		// "testHeader4": testHeader4,
		// "testHeader5": testHeader5,
		// "testHeader6": testHeader6,
	} {
		t.Run(scenario, f)
	}
}

func testHeader1(t *testing.T) {
	require.True(t, h1Regex.MatchString("# Header1"))
	// Returns true for 1 match
	require.False(t, h1Regex.MatchString("## Header2"), "## Header2")
	require.False(t, h1Regex.MatchString("            #Header1"))
	// This is what we can do, if we find a match.
	require.False(t, h1Regex.MatchString("gibberish"))
}

func testHeader2(t *testing.T) {
	require.True(t, h2Regex.MatchString("## Header1"))
	// This is what we can do, if we find a match.
	require.True(t, "## Header2" != h1Regex.FindString("# Header2"))
	require.True(t, "## Header2" != h1Regex.FindString("### Header2"))
	require.True(t, h2Regex.MatchString("## Header2"))
	require.False(t, h2Regex.MatchString("gibberish"))
}

func testHeader3(t *testing.T) {
	require.True(t, h3Regex.MatchString("# Header1"))
	require.False(t, h3Regex.MatchString("## Header2"))
	require.False(t, h3Regex.MatchString("gibberish"))
}

func testHeader4(t *testing.T) {
	require.True(t, h4Regex.MatchString("# Header1"))
	require.False(t, h4Regex.MatchString("## Header2"))
	require.False(t, h4Regex.MatchString("gibberish"))
}

func testHeader5(t *testing.T) {
	require.True(t, h5Regex.MatchString("# Header1"))
	require.False(t, h5Regex.MatchString("## Header2"))
	require.False(t, h5Regex.MatchString("gibberish"))
}

func testHeader6(t *testing.T) {
	require.True(t, h6Regex.MatchString("# Header1"))
	require.False(t, h6Regex.MatchString("## Header2"))
	require.False(t, h6Regex.MatchString("gibberish"))
}
