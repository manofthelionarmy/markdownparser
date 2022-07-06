package markdown

import (
	"fmt"
	"strings"
)

func convert(b string) []byte {
	var text string
	// need to refactor regex stuff
	if h1Regex.MatchString(b) {
		text = strings.Split(b, "#")[1]

		text = strings.TrimSpace(text)
		text = "<h1>" + text + "</h1>"
		text = parseNestedLink(text)
	} else if h2Regex.MatchString(b) {
		text = strings.Split(b, "##")[1]

		text = strings.TrimSpace(text)
		text = "<h2>" + text + "</h2>"
		text = parseNestedLink(text)
	} else if h3Regex.MatchString(b) {
		text = strings.Split(b, "###")[1]

		text = strings.TrimSpace(text)
		text = "<h3>" + text + "</h3>"
		text = parseNestedLink(text)
	} else if h4Regex.MatchString(b) {
		text = strings.Split(b, "####")[1]

		text = strings.TrimSpace(text)
		text = "<h4>" + text + "</h4>"
		text = parseNestedLink(text)
	} else if h5Regex.MatchString(b) {
		text = strings.Split(b, "#####")[1]

		text = strings.TrimSpace(text)
		text = "<h5>" + text + "</h5>"
		text = parseNestedLink(text)
	} else if h6Regex.MatchString(b) {
		text = strings.Split(b, "######")[1]

		text = strings.TrimSpace(text)
		text = "<h6>" + text + "</h6>"
		text = parseNestedLink(text)
	} else if strictLinkRegex.MatchString(b) {
		// Extract linkText by splitting by http text
		text = extractLink(b)
	} else if imageRegex.MatchString(b) {
		imgTag := "<img src=\"%s\" alt=\"%s\"/>"
		altText := imageAltRegex.FindString(b)
		src := imagePathRegex.FindString(b)
		text = fmt.Sprintf(imgTag, src[1:], strings.Trim(altText, "[]"))
	} else if unorderedListRegex.MatchString(b) {
		// I have an issue... what if we have something that is not a list underneath
		text = strings.ReplaceAll(b, "* ", "<li>")
		text = strings.ReplaceAll(text, "\n", "</li>\n")
		text = "<ul>\n" + text + "</ul>"
		// lists may have links!!!
		for hasNestedLink(text) {
			text = parseNestedLink(text)
		}
	} else {
		text = parseNestedLink(b)
		text = "<p>" + text + "</p>"
	}
	return []byte(text + "\n")
}

func hasNestedLink(b string) bool {
	return linkRegex.MatchString(b)
}

func extractLink(b string) string {
	var text string
	anchorText := linkTextRegex.FindString(b)
	link := httpRegex.FindString(b)
	altText := altTextRegex.FindString(b)
	if len(altText) > 0 {
		text = fmt.Sprintf("<a href=\"%s\" alt=%s>%s</a>", link, altText, anchorText[1:len(anchorText)-1])
	} else {
		text = fmt.Sprintf("<a href=\"%s\">%s</a>", link, anchorText[1:len(anchorText)-1])
	}
	return text
}

func parseNestedLink(b string) string {
	if hasNestedLink(b) {
		nestedLink := linkRegex.FindString(b)
		nestedLink = extractLink(nestedLink)
		splitedText := linkRegex.Split(b, 2)
		if len(splitedText) > 1 {
			b = splitedText[0] + nestedLink + splitedText[1]
		}
	}
	return b
}
