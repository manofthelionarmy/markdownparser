package markdown

import "regexp"

// These regexp work best if we trim left and right
var (
	headerOne   = `(?:^|\n)(#{1}\s)(.*)`
	headerTwo   = `(?:^|\n)(#{2}\s)(.*)`
	headerThree = `(?:^|\n)(#{3}\s)(.*)`
	headerFour  = `(?:^|\n)(#{4}\s)(.*)`
	headerFive  = `(?:^|\n)(#{5}\s)(.*)`
	headerSix   = `(?:^|\n)(#{6}\s)(.*)`
)

// TODO:
// 1. Links
// 2. Tables

// Headers
var (
	h1Regex *regexp.Regexp
	h2Regex *regexp.Regexp
	h3Regex *regexp.Regexp
	h4Regex *regexp.Regexp
	h5Regex *regexp.Regexp
	h6Regex *regexp.Regexp
)

func init() {
	h1Regex = regexp.MustCompile(headerOne)
	h2Regex = regexp.MustCompile(headerTwo)
	h3Regex = regexp.MustCompile(headerThree)
	h4Regex = regexp.MustCompile(headerFour)
	h5Regex = regexp.MustCompile(headerFive)
	h6Regex = regexp.MustCompile(headerSix)
}
