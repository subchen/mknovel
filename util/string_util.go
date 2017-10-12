package util

import (
	"strings"
)

func SubstrBefore(str string, find string) string {
	if len(find) == 0 {
		return ""
	}
	pos := strings.Index(str, find)
	if pos == -1 {
		return str
	}
	return str[:pos]
}

func SubstrAfter(str string, find string) string {
	if len(find) == 0 {
		return str
	}
	pos := strings.Index(str, find)
	if pos == -1 {
		return ""
	}
	return str[pos+len(find):]
}

func SubstrBeforeLast(str string, find string) string {
	if len(find) == 0 {
		return str
	}
	pos := strings.LastIndex(str, find)
	if pos == -1 {
		return str
	}
	return str[:pos]
}

func SubstrAfterLast(str string, find string) string {
	if len(find) == 0 {
		return ""
	}
	pos := strings.LastIndex(str, find)
	if pos == -1 {
		return ""
	}
	return str[pos+len(find):]
}

func SubstrBetween(str string, begin string, end string) string {
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
