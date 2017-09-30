package main

import (
	"fmt"
	"github.com/ungerik/go-dry"
	"hash/crc32"
	"html"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	RE_HTML_COMMENT = regexp.MustCompile("<!--.*-->")

	HTML_LINE_BREAK = []string{
		"<p>", "</p>", "<p/>", "<p />",
		"<br>", "<br/>", "<br />",
	}
	//HTML_WHITESPACES = []string{
	//	"&nbsp;", "�", "",
	//}

	NEW_LINE        = "\r\n"
	NEW_LINE_DOUBLE = NEW_LINE + NEW_LINE

	UNUSED_KEYWORDS = []string{
		"月票", "推荐票", "打赏",
		"书友", "兄弟姐妹",
		"微信", "公众号", "公告",
		"喜欢这部作品", "收藏", "下载链接", "注册", "您的支持",
		"未完待续", "请假一天", "新的一周",
		"第一更", "第二更", "第三更",
	}
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

func htmlAsTextLines(htmltext string, trimTrailingAd bool) []string {
	text := htmltext

	for _, tag := range HTML_LINE_BREAK {
		text = strings.Replace(text, tag, "\n", -1)
		text = strings.Replace(text, strings.ToUpper(tag), "\n", -1)
	}
	//for _, space := range HTML_WHITESPACES {
	//	text = strings.Replace(text, space, " ", -1)
	//}

	text = RE_HTML_COMMENT.ReplaceAllString(text, "")
	text = html.UnescapeString(text)

	var lines []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	if trimTrailingAd && len(lines) > 10 {
		// remove last unused line
		for i := 0; i < 3; i++ {
			line := lines[len(lines)-1]
			if !isUnusedLine(line) {
				break
			}
			// remove it, and continue
			lines = lines[:len(lines)-1]
		}
	}

	return lines
}

func isUnusedLine(text string) bool {
	for _, keyword := range UNUSED_KEYWORDS {
		if strings.Index(text, keyword) >= 0 {
			fmt.Println("  >> " + text)
			return true
		}
	}
	return false
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

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func fileGetAsString(file string, charset string) string {
	data, err := dry.FileGetBytes(file)
	panicIfError(err)

	return string(decodeBytes(data, charset))
}
