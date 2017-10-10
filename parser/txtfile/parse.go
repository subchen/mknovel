package txtfile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/util"
	"github.com/ungerik/go-dry"
)

func ImportAndParse(novel *model.Novel, inputEncoding string) {
	if novel.Name == "" {
		panic("no novel name provided")
	}
	if novel.Author == "" {
		panic("no novel author provided")
	}

	// download or copy cover-image to cache dir
	if novel.CoverImageURL != "" {
		fmt.Printf("Novel Cover Image: %v\n", novel.CoverImageURL)
		coverImageFile := filepath.Join(novel.CacheDirectory, "raw/cover.jpg")
		if !dry.FileExists(coverImageFile) {
			coverImageBytes, err := dry.FileGetBytes(novel.CoverImageURL)
			dry.PanicIfErr(err)
			err = dry.FileSetBytes(coverImageFile, coverImageBytes)
			dry.PanicIfErr(err)
		}
	}

	fmt.Println()

	fileBytes, err := dry.FileGetBytes(novel.URL)
	dry.PanicIfErr(err)

	txt := string(util.DecodeBytes(fileBytes, inputEncoding))
	txt = strings.Replace(txt, "\r", "", -1)

	txtChapterList := strings.Split(txt, "\n\n\n")
	for _, txtChapter := range txtChapterList {
		addChapter(novel, txtChapter)
	}
}

func addChapter(novel *model.Novel, txtChapter string) {
	var lines []string
	for _, line := range strings.Split(txtChapter, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	if len(lines) < 2 {
		return
	}

	index := len(novel.ChapterList) + 1
	chapter := &model.NovelChapter{
		ID:        fmt.Sprintf("%04d", index),
		Index:     index, // from base 1
		Name:      strings.TrimSpace(lines[0]),
		TextLines: lines[1:],
	}

	fmt.Printf("Found: %s %s\n", chapter.ID, chapter.Name)

	novel.ChapterList = append(novel.ChapterList, chapter)
}
