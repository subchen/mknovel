package main

import (
	"fmt"
	"github.com/ungerik/go-dry"
	"os"
	"path/filepath"
	"strings"
)

func packageNovelAsZIP(novel *Novel, outputDirectory string, txtEncoding string, zipFilenameEncoding string) {
	fmt.Println()
	fmt.Printf("Generating %d chapter files ...\n", len(novel.ChapterList))

	dir := filepath.Join(novel.CacheDirectory, "zip")
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	for _, chapter := range novel.ChapterList {
		txt := strings.Join(chapter.Lines, NEW_LINE_DOUBLE)
		txt = chapter.Name + NEW_LINE_DOUBLE + txt

		data := encodeBytes([]byte(txt), txtEncoding)

		file := filepath.Join(dir, chapter.FileId()+"-"+chapter.Name+".txt")
		dry.FileSetBytes(file, data)
	}

	// zip novel file
	file := filepath.Join(outputDirectory, fmt.Sprintf("%s-%s.zip", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Archiving Zip: %v ...\n", file)
	zipToFile(file, dir, zipFilenameEncoding)
}
