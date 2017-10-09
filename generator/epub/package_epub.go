package epub

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/util"
	"github.com/ungerik/go-dry"
	"net/url"
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
	srcFile := filepath.Join(util.GetCurrentDirectory(), "generator/epub", templateFile)

	dest, err := os.Create(destFile)
	dry.PanicIfErr(err)

	tmpl, err := template.New("").ParseFiles(srcFile)
	dry.PanicIfErr(err)

	err = tmpl.ExecuteTemplate(dest, filepath.Base(srcFile), context)
	dry.PanicIfErr(err)
}
