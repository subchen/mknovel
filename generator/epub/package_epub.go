package epub

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"text/template"

	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/util"
	"github.com/ungerik/go-dry"
)

var (
	templateCache = make(map[string]*template.Template)
)

func PackageNovelAsEPUB(novel *model.Novel, outputDirectory string) {
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
	executeTemplate("template/mimetype", filepath.Join(dir, "mimetype"), nil)
	executeTemplate("template/META-INF/container.xml", filepath.Join(dir, "META-INF/container.xml"), nil)
	executeTemplate("template/OEBPS/content.opf", filepath.Join(dir, "OEBPS/content.opf"), novel)
	executeTemplate("template/OEBPS/toc.ncx", filepath.Join(dir, "OEBPS/toc.ncx"), novel)
	executeTemplate("template/OEBPS/css/style.css", filepath.Join(dir, "OEBPS/css/style.css"), nil)
	executeTemplate("template/OEBPS/data/copyrights.xhtml", filepath.Join(dir, "OEBPS/data/copyrights.xhtml"), novel)

	for _, chapter := range novel.ChapterList {
		destFile := filepath.Join(dir, "OEBPS/data", chapter.ID+".xhtml")
		executeTemplate("template/OEBPS/data/chapter.xhtml", destFile, chapter)
	}

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

func executeTemplate(templateFile string, destFile string, context interface{}) {
	fmt.Printf("Writing: %v ...\n", destFile)

	// cache template obj
	t, ok := templateCache[templateFile]
	if !ok {
		var err error
		srcFileBytes := MustAsset(templateFile)
		t, err = template.New("").Parse(string(srcFileBytes))
		dry.PanicIfErr(err)

		templateCache[templateFile] = t
	}

	dest, err := os.Create(destFile)
	dry.PanicIfErr(err)
	defer dest.Close()

	err = t.Execute(dest, context)
	dry.PanicIfErr(err)
}
