package txt

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/subchen/go-stack"
	"github.com/subchen/mknovel/model"
	"github.com/ungerik/go-dry"
)

var (
	SINGLE_NEW_LINE_BYTES = []byte("\r\n")
	DOUBLE_NEW_LINE_BYTES = []byte("\r\n\r\n")
)

func PackageNovelAsTXT(novel *model.Novel, outputDirectory string, outputEncoding string, isDebug bool) {
	file := filepath.Join(outputDirectory, fmt.Sprintf("%s-%s.txt", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Writing TXT: %v ...\n", file)

	w, err := os.Create(file)
	dry.PanicIfErr(err)
	defer w.Close()

	writeln(w, "书名："+novel.Name, outputEncoding)
	writeln(w, "作者："+novel.Author, outputEncoding)

	for _, chapter := range novel.ChapterList {
		w.Write(DOUBLE_NEW_LINE_BYTES)
		writeln(w, chapter.Name, outputEncoding)
		for _, line := range chapter.TextLines {
			writeln(w, "　　"+line, outputEncoding)
		}
	}
}

func writeln(w io.Writer, s string, outputEncoding string) {
	if s != "" {
		data := gstack.CharsetEncodeBytes([]byte(s), outputEncoding)
		w.Write(data)
	}
	w.Write(SINGLE_NEW_LINE_BYTES)
}
