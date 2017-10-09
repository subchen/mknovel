package main

import (
	"fmt"
	"github.com/ungerik/go-dry"
	"os"
	"path/filepath"
	"text/template"
)

func packageNovelAsEPUB(novel *Novel, outputDirectory string) {
	fmt.Println()
	fmt.Printf("Generating %d chapter files ...\n", len(novel.ChapterList))

	dir := filepath.Join(novel.CacheDirectory, "epub")
	os.RemoveAll(dir)
	os.MkdirAll(filepath.Join(dir, "META-INF"), 0755)
	os.MkdirAll(filepath.Join(dir, "OEBPS/images"), 0755)
	os.MkdirAll(filepath.Join(dir, "OEBPS/css"), 0755)
	os.MkdirAll(filepath.Join(dir, "OEBPS/data"), 0755)

	// generate files
	executeTemplate("templates/epub_v2/mimetype", filepath.Join(dir, "mimetype"), nil)
	executeTemplate("templates/epub_v2/META-INF/container.xml", filepath.Join(dir, "META-INF/container.xml"), nil)
	executeTemplate("templates/epub_v2/OEBPS/content.opf", filepath.Join(dir, "OEBPS/content.opf"), novel)
	executeTemplate("templates/epub_v2/OEBPS/toc.ncx", filepath.Join(dir, "OEBPS/toc.ncx"), novel)
	executeTemplate("templates/epub_v2/OEBPS/css/style.css", filepath.Join(dir, "OEBPS/css/style.css"), nil)
	executeTemplate("templates/epub_v2/OEBPS/data/copyrights.xhtml", filepath.Join(dir, "OEBPS/data/copyrights.xhtml"), novel)

	for _, chapter := range novel.ChapterList {
		destFile := filepath.Join(dir, "OEBPS/data", chapter.FileId()+".xhtml")
		executeTemplate("templates/epub_v2/OEBPS/data/chapter.xhtml", destFile, chapter)
	}

	// copy cover-image
	if novel.CoverImageLink != "" {
		coverImageSrc := filepath.Join(novel.CacheDirectory, filepath.Base(novel.CoverImageLink))
		dry.FileCopy(coverImageSrc, filepath.Join(dir, "OEBPS/images/cover.jpg"))
	}

	// zip novel file
	file := filepath.Join(outputDirectory, fmt.Sprintf("%s-%s.epub", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Archiving ePub: %v ...\n", file)
	zipToFile(file, dir, "UTF-8")
}

func executeTemplate(templateFile string, destFile string, context interface{}) {
	srcFile := filepath.Join(getCurrentDirectory(), templateFile)

	dest, err := os.Create(destFile)
	panicIfError(err)

	tmpl, err := template.New("").ParseFiles(srcFile)
	panicIfError(err)

	err = tmpl.ExecuteTemplate(dest, filepath.Base(srcFile), context)
	panicIfError(err)
}
