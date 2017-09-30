package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type (
	Novel struct {
		BookURL        *url.URL
		Config         *NovelConfig
		CacheDirectory string

		Name           string
		Author         string
		Description    string
		CoverImageLink string
		ChapterList    []*NovelChapter

		BookId      string
		Publisher   string
		PublishDate string
	}

	NovelChapter struct {
		Index int
		Name  string
		Link  string

		DiskFile string
		Lines    []string
		Size     int
	}
)

func newNovel(bookUrl *url.URL, outputDirectory string) *Novel {
	novel := &Novel{
		BookURL:     bookUrl,
		Config:      loadConfigFile(bookUrl.Host + ".yaml"),
		BookId:      hashCRC(bookUrl.String()),
		Publisher:   "https://github.com/subchen/mknovel",
		PublishDate: time.Now().String(),
	}
	novel.CacheDirectory = filepath.Join(outputDirectory, ".cache", novel.BookId)

	// make cache dir
	os.MkdirAll(novel.CacheDirectory, 0755)
	fmt.Printf("Novel URL: %s\n", novel.BookURL)
	fmt.Printf("Output cache directory: %s\n", novel.CacheDirectory)

	return novel
}

func getNovelTitle(html string, config *NovelTitleConfig) string {
	if config.Begin != "" && config.End != "" {
		html = substrBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindStringSubmatch(html)
	title := matches[config.NameIndex]
	return strings.TrimSpace(title)
}

func getNovelAuthor(html string, config *NovelAuthorConfig) string {
	if config.Begin != "" && config.End != "" {
		html = substrBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindStringSubmatch(html)
	title := matches[config.NameIndex]
	return strings.TrimSpace(title)
}

func getNovelCoverImage(html string, config *NovelCoverImageConfig) string {
	if config.Begin != "" && config.End != "" {
		html = substrBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindStringSubmatch(html)
	title := matches[config.NameIndex]
	return strings.TrimSpace(title)
}

func getNovelChapterList(html string, config *NovelChapterIndexConfig) []*NovelChapter {
	if config.Begin != "" && config.End != "" {
		html = substrBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindAllStringSubmatch(html, -1)
	chapterList := make([]*NovelChapter, 0, len(matches))

	for i, match := range matches {
		chapter := &NovelChapter{
			Index: i + 1, // from base 1
			Name:  strings.TrimSpace(match[config.NameIndex]),
			Link:  match[config.LinkIndex],
		}
		chapterList = append(chapterList, chapter)
	}
	return chapterList
}

func getNovelChapterContent(html string, config *NovelChapterContentConfig) string {
	html = substrBetween(html, config.Begin, config.End)
	return html
}

func (c *NovelChapter) FileId() string {
	return fmt.Sprintf("%04d", c.Index)
}
