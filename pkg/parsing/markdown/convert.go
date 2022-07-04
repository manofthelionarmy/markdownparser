package markdown

import (
	"strings"
)

func convert(b string) []byte {
	var text string
	if h1Regex.MatchString(b) {
		text = strings.Split(b, "#")[1]

		text = strings.TrimSpace(text)
		text = "<h1>" + text + "</h1>"

		return []byte(text + "\n")
	} else if h2Regex.MatchString(b) {
		text = strings.Split(b, "##")[1]

		text = strings.TrimSpace(text)
		text = "<h2>" + text + "</h2>"
	} else if h3Regex.MatchString(b) {
		text = strings.Split(b, "###")[1]

		text = strings.TrimSpace(text)
		text = "<h3>" + text + "</h3>"
	} else if h4Regex.MatchString(b) {
		text = strings.Split(b, "####")[1]

		text = strings.TrimSpace(text)
		text = "<h4>" + text + "</h4>"
	} else if h5Regex.MatchString(b) {
		text = strings.Split(b, "#####")[1]

		text = strings.TrimSpace(text)
		text = "<h5>" + text + "</h5>"
	} else if h6Regex.MatchString(b) {
		text = strings.Split(b, "######")[1]

		text = strings.TrimSpace(text)
		text = "<h6>" + text + "</h6>"
	} else {
		text = "<p>" + b + "</p>"
	}
	return []byte(text + "\n")
}
