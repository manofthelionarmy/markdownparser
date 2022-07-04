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
		text = linkTextRegex.FindString(b)
		linkAndAlt := httpRegex.FindString(b)
		splits := strings.Split(linkAndAlt, "\"")
		link := splits[0]
		link = strings.TrimSpace(link)
		if len(splits) == 1 {
			// weird edge case
			text = fmt.Sprintf("<a href=\"%s\">%s</a>", link[:len(link)-1], text[1:len(text)-1])
		} else {
			alt := splits[1]
			text = fmt.Sprintf("<a href=\"%s\" alt=\"%s\">%s</a>", link, alt, text[1:len(text)-1])
		}
	} else if imageRegex.MatchString(b) {
		imgTag := "<img src=\"%s\" alt=\"%s\"/>"
		altText := imageAltRegex.FindString(b)
		src := imagePathRegex.FindString(b)
		text = fmt.Sprintf(imgTag, src[1:], strings.Trim(altText, "[]"))
	} else {
		text = "<p>" + b + "</p>"
		text = parseNestedLink(text)
	}
	return []byte(text + "\n")
}

func hasNestedLink(b string) bool {
	return linkRegex.MatchString(b)
}

func extractLink(b string) string {
	text := linkTextRegex.FindString(b)
	linkAndAlt := httpRegex.FindString(b)
	splits := strings.Split(linkAndAlt, "\"")
	link := splits[0]
	link = strings.TrimSpace(link)
	if len(splits) == 1 {
		// weird edge case
		text = fmt.Sprintf("<a href=\"%s\">%s</a>", link[:len(link)-1], text[1:len(text)-1])
	} else {
		alt := splits[1]
		text = fmt.Sprintf("<a href=\"%s\" alt=\"%s\">%s</a>", link, alt, text[1:len(text)-1])
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
