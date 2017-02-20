package main

import (
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func getEncoding(charset string) encoding.Encoding {
	name := strings.ToUpper(charset)

	if name == "GBK" || name == "GB2312" {
		return simplifiedchinese.GBK
	}
	if name == "GB18030" {
		return simplifiedchinese.GB18030
	}

	panic("Unsupport charset: " + charset)
}

func decodeString(text []byte, charset string) string {
	decoder := getEncoding(charset).NewDecoder()

	dst := make([]byte, len(text)*2)
	n, _, err := decoder.Transform(dst, text, true)
	if err != nil {
		panic(err)
	}
	return string(dst[:n])
}

func encodeString(text []byte, charset string) string {
	encoder := getEncoding(charset).NewEncoder()

	dst := make([]byte, len(text)*2)
	n, _, err := encoder.Transform(dst, text, true)
	if err != nil {
		panic(err)
	}
	return string(dst[:n])
}
