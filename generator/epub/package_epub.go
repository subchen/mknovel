package epub

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/util"
	"github.com/ungerik/go-dry"
	"github.com/wushilin/threads"
)

func PackageNovelAsEPUB(novel *model.Novel, outputDirectory string, isDebug bool) {
	fmt.Println()
	fmt.Printf("Generating %d chapter files ...\n", len(novel.ChapterList))

	dir := filepath.Join(novel.CacheDirectory, "epub")
	os.RemoveAll(dir)
	os.MkdirAll(filepath.Join(dir, "META-INF"), 0755)
	os.MkdirAll(filepath.Join(dir, "OEBPS/images"), 0755)
	os.MkdirAll(filepath.Join(dir, "OEBPS/css"), 0755)
	os.MkdirAll(filepath.Join(dir, "OEBPS/data"), 0755)

	// using host as source
	if webURL, err := url.Parse(novel.URL); err == nil {
		if webURL.Host != "" {
			novel.Source = webURL.Host
		}
	}

	// generate files
	executeTemplate("template/mimetype", filepath.Join(dir, "mimetype"), nil, isDebug)
	executeTemplate("template/META-INF/container.xml", filepath.Join(dir, "META-INF/container.xml"), nil, isDebug)
	executeTemplate("template/OEBPS/content.opf", filepath.Join(dir, "OEBPS/content.opf"), novel, isDebug)
	executeTemplate("template/OEBPS/toc.ncx", filepath.Join(dir, "OEBPS/toc.ncx"), novel, isDebug)
	executeTemplate("template/OEBPS/css/style.css", filepath.Join(dir, "OEBPS/css/style.css"), nil, isDebug)
	executeTemplate("template/OEBPS/data/copyrights.xhtml", filepath.Join(dir, "OEBPS/data/copyrights.xhtml"), novel, isDebug)

	// create chapters using thread-pool
	nThreads := 100
	pool := threads.NewPool(nThreads, nThreads*2)
	poolMutux := &sync.Mutex{}
	pool.Start()
	for _, chapter := range novel.ChapterList {
		pool.Submit(writeChapterFile(chapter, dir, poolMutux, isDebug))
	}
	pool.Shutdown()
	pool.Wait()

	// copy cover-image
	if novel.CoverImageURL != "" {
		coverImageSrc := filepath.Join(novel.CacheDirectory, "raw/cover.jpg")
		dry.FileCopy(coverImageSrc, filepath.Join(dir, "OEBPS/images/cover.jpg"))
	}

	// zip novel file
	file := filepath.Join(outputDirectory, fmt.Sprintf("%s-%s.epub", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Archiving ePub: %v ...\n", file)
	util.ZipToFile(file, dir, "UTF-8")
}

func writeChapterFile(chapter *model.NovelChapter, dir string, poolMutux *sync.Mutex, isDebug bool) threads.JobFunc {
	return func() interface{} {
		destFile := filepath.Join(dir, "OEBPS/data", chapter.ID+".xhtml")
		if isDebug {
			poolMutux.Lock()
			fmt.Printf("Writing %s ...\n", destFile)
			poolMutux.Unlock()
		}
		executeTemplate("template/OEBPS/data/chapter.xhtml", destFile, chapter, isDebug)
		return nil
	}
}

func executeTemplate(templateFile string, destFile string, context interface{}, isDebug bool) {
	srcFileBytes := MustAsset(templateFile)
	t, err := template.New("").Parse(string(srcFileBytes))
	dry.PanicIfErr(err)

	dest, err := os.Create(destFile)
	dry.PanicIfErr(err)
	defer dest.Close()

	err = t.Execute(dest, context)
	dry.PanicIfErr(err)
}
