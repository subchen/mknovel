package txtfile

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/subchen/go-stack"
	"github.com/subchen/mknovel/model"
	"github.com/ungerik/go-dry"
)

var (
	REGEXP_WHITESPACES_LINE = regexp.MustCompile(`(?m)^(\s|　)*$`)
)

func ImportAndParse(novel *model.Novel, inputEncoding string, autoChapterGroup bool) {
	fileBytes, err := dry.FileGetBytes(novel.URL)
	dry.PanicIfErr(err)

	txt := string(gstack.CharsetDecodeBytes(fileBytes, inputEncoding))
	txt = strings.Replace(txt, "\r", "", -1)
	txt = REGEXP_WHITESPACES_LINE.ReplaceAllString(txt, "\n")

	groupName := ""
	txtChapterList := strings.Split(txt, "\n\n\n")

	// try to get name and author from headers
	skipHeaderLines := 0
	headers := strings.Split(strings.TrimSpace(txtChapterList[0]), "\n")
	if len(headers) >= 2 {
		HEADER_KEYWORD := "作者"
		if strings.Contains(headers[0], HEADER_KEYWORD) {
			name := gstack.StringBefore(headers[0], HEADER_KEYWORD)
			author := gstack.StringAfter(headers[0], HEADER_KEYWORD)
			setNovelNameAndAuthor(novel, name, author)
			skipHeaderLines = 1
		} else if strings.Contains(headers[1], HEADER_KEYWORD) {
			setNovelNameAndAuthor(novel, headers[0], headers[1])
			skipHeaderLines = 2
		}
	}
	if skipHeaderLines > 0 {
		// reset frist Chapter
		txtChapterList[0] = strings.Join(headers[skipHeaderLines:], "\n")
	}

	// validate
	if novel.Name == "" {
		panic("no novel name provided")
	}
	if novel.Author == "" {
		panic("no novel author provided")
	}
	fmt.Println()
	fmt.Printf("Novel Name: %v\n", novel.Name)
	fmt.Printf("Novel Author: %v\n", novel.Author)

	fmt.Println()
	fmt.Println("Parsing chapter list ...")
	for _, txtChapter := range txtChapterList {
		addChapter(novel, txtChapter, &groupName, autoChapterGroup)
	}

	// download or copy cover-image to cache dir
	if novel.CoverImageURL != "" {
		fmt.Println()
		fmt.Printf("Novel Cover Image: %v\n", novel.CoverImageURL)
		coverImageFile := filepath.Join(novel.CacheDirectory, "raw/cover.jpg")
		if !dry.FileExists(coverImageFile) {
			coverImageBytes, err := dry.FileGetBytes(novel.CoverImageURL)
			dry.PanicIfErr(err)
			err = dry.FileSetBytes(coverImageFile, coverImageBytes)
			dry.PanicIfErr(err)
		}
	}
}

func addChapter(novel *model.Novel, txtChapter string, groupName *string, autoChapterGroup bool) {
	var lines []string
	for _, line := range strings.Split(txtChapter, "\n") {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return
	}

	if len(lines) == 1 {
		*groupName = lines[0]
		return
	}

	if autoChapterGroup && *groupName != "" {
		lines[0] = *groupName + " " + lines[0]
	}

	name := lines[0]
	name = strings.TrimSpace(name)
	name = strings.TrimSuffix(name, "：")
	name = strings.TrimSuffix(name, ":")

	index := len(novel.ChapterList) + 1
	chapter := &model.NovelChapter{
		ID:        fmt.Sprintf("%04d", index),
		Index:     index, // from base 1
		Name:      name,
		TextLines: lines[1:],
	}

	fmt.Printf("Found: %s %s\n", chapter.ID, chapter.Name)

	novel.ChapterList = append(novel.ChapterList, chapter)
}

func setNovelNameAndAuthor(novel *model.Novel, name string, author string) {
	if novel.Name == "" {
		name = strings.Replace(name, "书名", "", -1)
		name = strings.Replace(name, "：", "", -1)
		name = strings.Replace(name, ":", "", -1)
		name = strings.Replace(name, "《", "", -1)
		name = strings.Replace(name, "》", "", -1)
		novel.Name = strings.TrimSpace(name)
	}

	if novel.Author == "" {
		author = strings.Replace(author, "作者", "", -1)
		author = strings.Replace(author, "：", "", -1)
		author = strings.Replace(author, ":", "", -1)
		novel.Author = strings.TrimSpace(author)
	}
}
