package main

import (
	"hash/crc32"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	RE_HTML_COMMENT = regexp.MustCompile("<!--.*-->")
	RE_HTML_SPACE   = regexp.MustCompile("&nbsp;")
	RE_HTML_BR      = regexp.MustCompile("<br ?/?>")
	RE_HTML_P       = regexp.MustCompile("<p ?/?>")
)

func substrBetween(str string, begin string, end string) string {
	ipos := strings.Index(str, begin)
	if ipos < 0 {
		return ""
	}

	leftStr := str[ipos+len(begin) : len(str)]

	jpos := strings.Index(leftStr, end)
	if jpos < 0 {
		return ""
	}

	return leftStr[0:jpos]
}

func htmlAsText(html string) string {
	text := html
	text = strings.Replace(text, "</p>", "", -1)
	text = RE_HTML_COMMENT.ReplaceAllString(text, "")
	text = RE_HTML_BR.ReplaceAllString(text, "\r\n")
	text = RE_HTML_P.ReplaceAllString(text, "\r\n\r\n")

	text = strings.TrimSpace(text)
	text = RE_HTML_SPACE.ReplaceAllString(text, " ")
	return text
}

func hashCRC(text string) string {
	n := crc32.ChecksumIEEE([]byte(text))
	return strconv.FormatUint(uint64(n), 10)
}

func getExecutorDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return dir
}

func getCurrentDirectory() string {
	dir := os.Getenv("PWD")
	if dir == "" {
		dir = "."
	}
	return dir
}

func fileExist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

func fileSize(file string) int {
	if f, err := os.Stat(file); err == nil {
		return int(f.Size())
	}
	return 0
}
