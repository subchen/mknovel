package model

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ungerik/go-dry"
)

type (
	Novel struct {
		// basic
		Name          string
		Author        string
		CoverImageURL string

		// chapters
		ChapterList []*NovelChapter

		// for epub
		ID          string
		URL         string
		Source      string
		Subject     string
		Publisher   string
		PublishDate string

		// cache
		CacheDirectory string
	}

	NovelChapter struct {
		ID    string // == format("%04d", index)
		Index int    // index >= 1
		Name  string
		URL   string

		TextLines []string
		Size      int
	}
)

func NewNovel(opts *NovelOptions) *Novel {
	id := dry.StringMD5Hex(opts.NovelURL)
	outputDirectory, _ := filepath.Abs(opts.OutputDirectory)

	novel := &Novel{
		ID:          id,
		URL:         opts.NovelURL,
		Subject:     "mknovel",
		Publisher:   "https://github.com/subchen/mknovel",
		PublishDate: time.Now().String(),

		Name:          opts.NovelName,
		Author:        opts.NovelAuthor,
		CoverImageURL: opts.NovelCoverImageURL,

		CacheDirectory: filepath.Join(outputDirectory, ".cache", id),
	}

	// make cache dir
	os.MkdirAll(filepath.Join(novel.CacheDirectory, "raw"), 0755)
	fmt.Printf("Novel URL: %s\n", novel.URL)
	fmt.Printf("Output Cache Directory: %s\n", novel.CacheDirectory)

	return novel
}
