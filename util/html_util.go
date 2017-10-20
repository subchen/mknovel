package util

import (
	"html"
	"regexp"
	"strings"
)

var (
	RE_HTML_COMMENT = regexp.MustCompile("<!--.*-->")

	HTML_LINE_BREAK = []string{
		"<p>", "</p>", "<p/>", "<p />",
		"<br>", "<br/>", "<br />",
	}
)

func HtmlAsTextLines(htmltext string) []string {
	text := htmltext

	for _, tag := range HTML_LINE_BREAK {
		text = strings.Replace(text, tag, "\n", -1)
		text = strings.Replace(text, strings.ToUpper(tag), "\n", -1)
	}

	text = RE_HTML_COMMENT.ReplaceAllString(text, "")
	text = html.UnescapeString(text)

	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	return lines
}
