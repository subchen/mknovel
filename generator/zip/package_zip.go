package zip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/util"
	"github.com/ungerik/go-dry"
	"github.com/wushilin/threads"
)

const (
	NEW_LINES = "\r\n\r\n"
)

func PackageNovelAsZIP(novel *model.Novel, outputDirectory string, txtEncoding string, zipFilenameEncoding string, isDebug bool) {
	fmt.Println()
	fmt.Printf("Generating %d chapter files ...\n", len(novel.ChapterList))

	dir := filepath.Join(novel.CacheDirectory, "zip")
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	// create chapters using thread-pool
	nThreads := 100
	pool := threads.NewPool(nThreads, nThreads*2)
	poolMutux := &sync.Mutex{}
	pool.Start()
	for _, chapter := range novel.ChapterList {
		pool.Submit(writeChapterFile(chapter, dir, txtEncoding, poolMutux, isDebug))
	}
	pool.Shutdown()
	pool.Wait()

	// zip novel file
	file := filepath.Join(outputDirectory, fmt.Sprintf("%s-%s.zip", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Archiving Zip: %v ...\n", file)
	util.ZipToFile(file, dir, zipFilenameEncoding)
}

func writeChapterFile(chapter *model.NovelChapter, dir string, txtEncoding string, poolMutux *sync.Mutex, isDebug bool) threads.JobFunc {
	return func() interface{} {
		txt := chapter.Name + NEW_LINES + strings.Join(chapter.TextLines, NEW_LINES)
		data := util.EncodeBytes([]byte(txt), txtEncoding)

		destFile := filepath.Join(dir, chapter.ID+"-"+chapter.Name+".txt")
		if isDebug {
			poolMutux.Lock()
			fmt.Printf("Writing %s ...\n", destFile)
			poolMutux.Unlock()
		}
		dry.FileSetBytes(destFile, data)

		return nil
	}
}
