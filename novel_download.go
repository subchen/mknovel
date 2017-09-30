package main

import (
	"fmt"
	"github.com/ungerik/go-dry"
	"github.com/wushilin/threads"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func downloadNovel(novel *Novel, nThreads int, trimTrailingAd bool) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println()
			fmt.Printf("ERROR: %+v\n", err)
			fmt.Println()
			fmt.Println("You can retry later if hit network issue.")
			fmt.Println()
			os.Exit(1)
		}
	}()

	// download indexes
	fmt.Println()
	fmt.Printf("Downloading %s ...\n", novel.BookURL)
	html := downloadURLAsString(novel.BookURL.String(), novel.Config.WebsiteCharset)

	// parse title, author
	novel.Name = getNovelTitle(html, novel.Config.Title)
	fmt.Printf("Novel Name: %v\n", novel.Name)

	novel.Author = getNovelAuthor(html, novel.Config.Author)
	fmt.Printf("Novel Author: %v\n", novel.Author)

	// parse cover-image
	if novel.Config.CoverImage != nil {
		novel.CoverImageLink = getNovelCoverImage(html, novel.Config.CoverImage)
		fmt.Printf("Novel Cover Image: %v\n", novel.CoverImageLink)

		// download cover-image
		if strings.Contains(novel.CoverImageLink, "nocover") {
			novel.CoverImageLink = ""
		} else {
			coverImageURL := createURL(novel, novel.CoverImageLink)
			coverImageFile := filepath.Join(novel.CacheDirectory, filepath.Base(novel.CoverImageLink))
			downloadURLAsFile(coverImageURL, coverImageFile)
		}
	}

	// parse chapter list
	novel.ChapterList = getNovelChapterList(html, novel.Config.ChapterIndex)
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
		pool.Submit(downloadNovelChapter(novel, chapter, trimTrailingAd, downMutex, &downIndex))
	}
	pool.Shutdown()
	pool.Wait()
}

func downloadNovelChapter(novel *Novel, chapter *NovelChapter, trimTrailingAd bool, downMutex *sync.Mutex, downIndex *int) threads.JobFunc {
	return func() interface{} {
		// make chapter url
		chapterURL := createURL(novel, chapter.Link)

		downMutex.Lock()
		*downIndex = *downIndex + 1
		fmt.Printf("Downloading %d/%d %v ...\n", *downIndex, len(novel.ChapterList), chapterURL)
		downMutex.Unlock()

		// download chapter
		chapter.DiskFile = filepath.Join(novel.CacheDirectory, filepath.Base(chapterURL))
		downloadURLAsFile(chapterURL, chapter.DiskFile)

		// parse content
		html := fileGetAsString(chapter.DiskFile, novel.Config.WebsiteCharset)
		html = getNovelChapterContent(html, novel.Config.ChapterContent)

		// convert to plain txt lines
		chapter.Lines = htmlAsTextLines(html, trimTrailingAd)
		chapter.Size = len(strings.Join(chapter.Lines, "\n"))

		return nil
	}
}

func filterNovelShortChapters(novel *Novel, nShortChapter int) {
	var reCharpterList []*NovelChapter
	skippedShortChapter := 0
	for _, chapter := range novel.ChapterList {
		chapter.Index = len(reCharpterList) + 1 // reIndex

		if chapter.Size > nShortChapter {
			reCharpterList = append(reCharpterList, chapter)
		} else {
			if skippedShortChapter == 0 {
				fmt.Println()
				fmt.Println("Found Short Chapters (Skipped) ...")
			}
			skippedShortChapter++
			chapterURL := createURL(novel, chapter.Link)
			fmt.Printf("%v\n  %s (%d)\n", chapterURL, chapter.Name, chapter.Size)
		}
	}
	novel.ChapterList = reCharpterList
}

func downloadURLAsFile(url string, savedFile string) {
	// cache
	if dry.FileExists(savedFile) {
		return
	}

	res, err := http.Get(url)
	panicIfError(err)

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	panicIfError(err)

	err = dry.FileSetBytes(savedFile, data)
	panicIfError(err)
}

func downloadURLAsString(url string, charset string) string {
	res, err := http.Get(url)
	panicIfError(err)

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	panicIfError(err)

	return string(decodeBytes(data, charset))
}

func createURL(novel *Novel, link string) string {
	url, err := novel.BookURL.Parse(link)
	panicIfError(err)

	return url.String()
}
