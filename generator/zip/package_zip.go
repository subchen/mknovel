package zip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/util"
	"github.com/ungerik/go-dry"
)

const (
	NEW_LINES = "\r\n\r\n"
)

func PackageNovelAsZIP(novel *model.Novel, outputDirectory string, txtEncoding string, zipFilenameEncoding string) {
	fmt.Println()
	fmt.Printf("Generating %d chapter files ...\n", len(novel.ChapterList))

	dir := filepath.Join(novel.CacheDirectory, "zip")
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	for _, chapter := range novel.ChapterList {
		txt := chapter.Name + NEW_LINES + strings.Join(chapter.TextLines, NEW_LINES)
		bytes := util.EncodeBytes([]byte(txt), txtEncoding)

		file := filepath.Join(dir, chapter.ID+"-"+chapter.Name+".txt")
		dry.FileSetBytes(file, bytes)
	}

	// zip novel file
	file := filepath.Join(outputDirectory, fmt.Sprintf("%s-%s.zip", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Archiving Zip: %v ...\n", file)
	util.ZipToFile(file, dir, zipFilenameEncoding)
}
