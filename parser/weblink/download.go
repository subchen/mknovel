package weblink

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/subchen/go-stack/encoding/runes"
	"github.com/subchen/mknovel/model"
	"github.com/ungerik/go-dry"
	"github.com/wushilin/threads"
)

func StartDownload(novel *model.Novel, nThreads int, nShortChapterSize int) {
	webURL, err := url.Parse(novel.URL)
	dry.PanicIfErr(err)

	// load config
	config := loadConfig(webURL.Host + ".yaml")

	// download indexes
	fmt.Println()
	fmt.Printf("Downloading %s ...\n", novel.URL)
	htmlBytes, err := dry.FileGetBytes(novel.URL)
	dry.PanicIfErr(err)
	html := string(runes.DecodeBytes(htmlBytes, config.WebsiteCharset))

	// parse name
	if novel.Name == "" {
		novel.Name = getNovelTitle(html, config.Title)
		fmt.Printf("Novel Name: %v\n", novel.Name)
	}

	// parse author
	if novel.Author == "" {
		novel.Author = getNovelAuthor(html, config.Author)
		fmt.Printf("Novel Author: %v\n", novel.Author)
	}

	// parse cover-image
	if novel.CoverImageURL == "" && config.CoverImage != nil {
		coverImageLink := getNovelCoverImage(html, config.CoverImage)
		coverImageURL, err := webURL.Parse(coverImageLink)
		dry.PanicIfErr(err)

		if !strings.Contains(coverImageURL.String(), "nocover") {
			novel.CoverImageURL = coverImageURL.String()
			fmt.Printf("Novel Cover Image: %v\n", novel.CoverImageURL)
		}
	}
	// download or copy cover-image to cache dir
	if novel.CoverImageURL != "" {
		coverImageFile := filepath.Join(novel.CacheDirectory, "raw/cover.jpg")
		if !dry.FileExists(coverImageFile) {
			coverImageBytes, err := dry.FileGetBytes(novel.CoverImageURL)
			dry.PanicIfErr(err)
			err = dry.FileSetBytes(coverImageFile, coverImageBytes)
			dry.PanicIfErr(err)
		}
	}

	// parse chapter list
	novel.ChapterList = getNovelChapterList(html, config.ChapterIndex)
	fmt.Printf("Novel Chapter Count: %v\n", len(novel.ChapterList))

	if len(novel.ChapterList) < 3 {
		panic("Unable to parse chapter list")
	}

	fmt.Println()

	// download chapters using thread-pool
	pool := threads.NewPool(nThreads, nThreads*2)
	pool.Start()
	downMutex := &sync.Mutex{}
	downIndex := 0
	for _, chapter := range novel.ChapterList {
		// to abs url
		chapterURL, err := webURL.Parse(chapter.URL)
		dry.PanicIfErr(err)
		chapter.URL = chapterURL.String()

		// fork to download
		pool.Submit(downloadNovelChapter(novel, chapter, config, downMutex, &downIndex))
	}
	pool.Shutdown()
	pool.Wait()

	// remove short chapters
	skipNovelShortChapters(novel, nShortChapterSize)
}

func downloadNovelChapter(novel *model.Novel, chapter *model.NovelChapter, config *SiteConfig, downMutex *sync.Mutex, downIndex *int) threads.JobFunc {
	return func() interface{} {
		downMutex.Lock()
		*downIndex = *downIndex + 1
		fmt.Printf("Downloading %d/%d %v ...\n", *downIndex, len(novel.ChapterList), chapter.URL)
		downMutex.Unlock()

		// download chapter
		var chapterBytes []byte
		var err error
		chapterFile := filepath.Join(novel.CacheDirectory, "raw", filepath.Base(chapter.URL))
		if dry.FileExists(chapterFile) {
			// get from cache
			chapterBytes, err = dry.FileGetBytes(chapterFile)
			dry.PanicIfErr(err)
		} else {
			for i := 0; i < 10; i++ {
				if chapterBytes, err = dry.FileGetBytes(chapter.URL, 10*time.Second); err == nil {
					break
				}
			}
			dry.PanicIfErr(err)

			err = dry.FileSetBytes(chapterFile, chapterBytes)
			dry.PanicIfErr(err)
		}

		// parse content
		html := string(runes.DecodeBytes(chapterBytes, config.WebsiteCharset))
		html = getNovelChapterContent(html, config.ChapterContent)

		// convert to plain txt lines
		chapter.TextLines = htmlAsTextLines(html)
		chapter.Size = len(strings.Join(chapter.TextLines, "\n"))

		return nil
	}
}

func skipNovelShortChapters(novel *model.Novel, nShortChapterSize int) {
	var reChapterList []*model.NovelChapter
	skipped := 0
	for _, chapter := range novel.ChapterList {
		chapter.Index = len(reChapterList) + 1 // reIndex
		chapter.ID = fmt.Sprintf("%04d", chapter.Index)

		if chapter.Size > nShortChapterSize {
			reChapterList = append(reChapterList, chapter)
		} else {
			if skipped == 0 {
				fmt.Println()
				fmt.Println("Found Short Chapters (Skipped) ...")
			}
			skipped++
			fmt.Printf("%v\n  %s (%d)\n", chapter.URL, chapter.Name, chapter.Size)
		}
	}
	novel.ChapterList = reChapterList
}
