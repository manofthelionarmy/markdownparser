package markdown

import "regexp"

// These regexp work best if we trim left and right
var (
	// Headers
	headerOne    = `(?:^|\n)(#{1}\s)(.*)`
	headerTwo    = `(?:^|\n)(#{2}\s)(.*)`
	headerThree  = `(?:^|\n)(#{3}\s)(.*)`
	headerFour   = `(?:^|\n)(#{4}\s)(.*)`
	headerFive   = `(?:^|\n)(#{5}\s)(.*)`
	headerSix    = `(?:^|\n)(#{6}\s)(.*)`
	strictLink   = `(?:^|\n)(\[.*\])(\((http)(?:s)?(\:\/\/).*\))`
	link         = `(\[.*\])(\((http)(?:s)?(\:\/\/).*\))`
	linkText     = `(\[.*\])`
	httpText     = `((http)(?:s)?(\:\/\/).*)`
	image        = `(?:^|\n)(\!)(\[(?:.*)?\])(\(.*(\.(jpg|png|gif|tiff|bmp))(?:(\s\"|\')(\w|\W|\d)+(\"|\'))?\))`
	imageAltText = `(\[(?:.*)?\])`
	imageFile    = `(\(.*(\.(jpg|png|gif|tiff|bmp))(?:(\s\"|\')(\w|\W|\d)+(\"|\'))?\))`
)

// TODO:

// 2. Images
// 3. Unordered List
// 4. Tables

// Headers
var (
	h1Regex *regexp.Regexp
	h2Regex *regexp.Regexp
	h3Regex *regexp.Regexp
	h4Regex *regexp.Regexp
	h5Regex *regexp.Regexp
	h6Regex *regexp.Regexp

	strictLinkRegex *regexp.Regexp
	linkRegex       *regexp.Regexp
	linkTextRegex   *regexp.Regexp
	httpRegex       *regexp.Regexp

	imageRegex     *regexp.Regexp
	imageAltRegex  *regexp.Regexp
	imagePathRegex *regexp.Regexp
)

func init() {
	h1Regex = regexp.MustCompile(headerOne)
	h2Regex = regexp.MustCompile(headerTwo)
	h3Regex = regexp.MustCompile(headerThree)
	h4Regex = regexp.MustCompile(headerFour)
	h5Regex = regexp.MustCompile(headerFive)
	h6Regex = regexp.MustCompile(headerSix)

	strictLinkRegex = regexp.MustCompile(strictLink)
	linkRegex = regexp.MustCompile(link)
	linkTextRegex = regexp.MustCompile(linkText)
	httpRegex = regexp.MustCompile(httpText)

	imageRegex = regexp.MustCompile(image)
	imageAltRegex = regexp.MustCompile(imageAltText)
	imagePathRegex = regexp.MustCompile(imageFile)
}
