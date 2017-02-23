package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"

	"github.com/go-yaml/yaml"
	"github.com/wushilin/threads"
)

type (
	Novel struct {
		Name        string
		Author      string
		BookUrl     *url.URL
		ChapterList []NovelChapter
		Config      *NovelConfig
		Tempdir     string
	}
	NovelChapter struct {
		Index int
		Name  string
		Link  string
	}

	NovelConfig struct {
		WebsiteCharset     string              `yaml:"website-charset"`
		ZipFilenameCharset string              `yaml:"zipfilename-charset"`
		Title              *NovelTitleConfig   `yaml:"title"`
		Author             *NovelAuthorConfig  `yaml:"author"`
		Chapter            *NovelChapterConfig `yaml:"chapter`
		Content            *NovelContentConfig `yaml:"content"`
	}
	NovelTitleConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
	}
	NovelAuthorConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
	}
	NovelChapterConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
		LinkIndex int    `yaml:"link-index"`
	}
	NovelContentConfig struct {
		Begin           string `yaml:"begin"`
		End             string `yaml:"end"`
		IgnoreShortText int    `yaml:"ignore-short-text"`
	}
)

func findConfigFile(file string) string {
	configFile, _ := filepath.Abs(file)
	if fileExist(configFile) {
		return configFile
	}

	configFile = filepath.Join(getExecutorDirectory(), file)
	if fileExist(configFile) {
		return configFile
	}

	configFile = filepath.Join("/etc/mknovel", file)
	if fileExist(configFile) {
		return configFile
	}

	panic("not found config file: " + file)
}

func loadConfigFile(file string) *NovelConfig {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	config := NovelConfig{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	return &config
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

func getNovelChapterList(html string, config *NovelChapterConfig) []NovelChapter {
	if config.Begin != "" && config.End != "" {
		html = substrBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindAllStringSubmatch(html, -1)
	chapterList := make([]NovelChapter, 0, len(matches))

	for i, match := range matches {
		chapter := NovelChapter{
			Index: i + 1,
			Name:  strings.TrimSpace(match[config.NameIndex]),
			Link:  match[config.LinkIndex],
		}
		chapterList = append(chapterList, chapter)
	}
	return chapterList
}

func getNovelChapterContent(html string, config *NovelContentConfig) string {
	html = substrBetween(html, config.Begin, config.End)
	return html
}

func downloadURL(url string, charset string) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if charset != "" {
		return decodeString(data, charset)
	} else {
		return string(data)
	}
}

func downloadNovelChapter(novel *Novel, chapter NovelChapter, nIndex *int32) threads.JobFunc {
	return func() interface{} {
		// make chapter url
		chapterUrl, err := novel.BookUrl.Parse(chapter.Link)
		if err != nil {
			panic(err)
		}

		atomic.AddInt32(nIndex, 1)
		fmt.Printf("Downloading %d/%d %04d %v ...\n", *nIndex, len(novel.ChapterList), chapter.Index, chapterUrl.String())

		// make file name
		file := filepath.Join(novel.Tempdir, fmt.Sprintf("%04d %v.txt", chapter.Index, chapter.Name))
		if fileExist(file) {
			//continue
			return nil
		}

		// download
		html := downloadURL(chapterUrl.String(), novel.Config.WebsiteCharset)

		// html to text
		html = getNovelChapterContent(html, novel.Config.Content)
		if len(html) < novel.Config.Content.IgnoreShortText {
			fmt.Printf("Ignored short chapter %04d %s\n", chapter.Index, chapter.Name)
			//continue
			return nil
		}

		text := "    " + chapter.Name + "\n\n" + htmlAsText(html)

		// write file
		err = ioutil.WriteFile(file, []byte(text), 0644)
		if err != nil {
			panic(err)
		}

		return nil
	}
}

func downloadNovel(bookUrl *url.URL, dir string, nThreads int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println()
			fmt.Printf("ERROR: %+v\n", err)
			os.Exit(1)
		}
	}()

	novel := &Novel{}
	novel.BookUrl = bookUrl

	novel.Tempdir = filepath.Join(dir, hashCRC(bookUrl.String()))
	os.MkdirAll(novel.Tempdir, 0755)

	fmt.Printf("Novel URL: %s\n", novel.BookUrl)
	fmt.Printf("Output directory: %s\n", novel.Tempdir)

	// load config
	configFile := findConfigFile(bookUrl.Host + ".yaml")
	fmt.Printf("Config File: %s\n", configFile)
	novel.Config = loadConfigFile(configFile)

	// download indexes
	fmt.Println()
	fmt.Printf("Downloading %s ...\n", novel.BookUrl)
	html := downloadURL(bookUrl.String(), novel.Config.WebsiteCharset)

	novel.Name = getNovelTitle(html, novel.Config.Title)
	fmt.Printf("Novel Name: %v\n", novel.Name)
	novel.Author = getNovelAuthor(html, novel.Config.Author)
	fmt.Printf("Novel Author: %v\n", novel.Author)

	novel.ChapterList = getNovelChapterList(html, novel.Config.Chapter)
	fmt.Printf("Novel Chapter Count: %v\n", len(novel.ChapterList))

	if len(novel.ChapterList) < 3 {
		panic("Unable to parse chapter list")
	}

	pool := threads.NewPool(nThreads, nThreads*2)
	pool.Start()

	// download chapters
	nIndex := int32(0)
	fmt.Println()
	for _, chapter := range novel.ChapterList {
		pool.Submit(downloadNovelChapter(novel, chapter, &nIndex))
	}

	pool.Shutdown()
	pool.Wait()

	// zip novel file
	file := filepath.Join(dir, fmt.Sprintf("%s (%s).zip", novel.Name, novel.Author))
	fmt.Println()
	fmt.Printf("Archive Zip: %v ...\n", file)
	zipToFile(file, novel.Tempdir, novel.Config.ZipFilenameCharset)

	// done
	fmt.Println()
	fmt.Println("Completed!")
}
