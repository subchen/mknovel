package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
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
		Subject     string
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
		Subject:     "mknovel",
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

func (c *NovelChapter) FileId() string {
	return fmt.Sprintf("%04d", c.Index)
}
