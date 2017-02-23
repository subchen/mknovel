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

	"github.com/go-yaml/yaml"
)

type (
	NovelChapter struct {
		Index int
		Name  string
		Link  string
	}

	NovelConfig struct {
		WebsiteCharset     string              `yaml:"website-charset"`
		ZipFilenameCharset string              `yaml:"zipfilename-charset"`
		Title              *NovelTitleConfig   `yaml:"title"`
		Chapter            *NovelChapterConfig `yaml:"chapter`
		Content            *NovelContentConfig `yaml:"content"`
	}
	NovelTitleConfig struct {
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

func downloadNovel(bookUrl *url.URL, dir string) {
	tempdir := filepath.Join(dir, hashCRC(bookUrl.String()))
	os.MkdirAll(tempdir, 0755)

	fmt.Printf("Novel URL: %s\n", bookUrl)
	fmt.Printf("Output directory: %s\n", tempdir)

	// load config
	configFile := findConfigFile(bookUrl.Host + ".yaml")
	fmt.Printf("Config File: %s\n", configFile)
	config := loadConfigFile(configFile)

	// download indexes
	fmt.Println()
	fmt.Printf("Downloading %s ...\n", bookUrl)
	html := downloadURL(bookUrl.String(), config.WebsiteCharset)

	title := getNovelTitle(html, config.Title)
	fmt.Printf("Novel Name: %v\n", title)

	chapterList := getNovelChapterList(html, config.Chapter)
	fmt.Printf("Novel Chapter Count: %v\n", len(chapterList))

	if len(chapterList) < 3 {
		panic("Unable to parse chapter list")
	}

	// download chapters
	fmt.Println()
	for _, chapter := range chapterList {
		// make chapter url
		chapterUrl, err := bookUrl.Parse(chapter.Link)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Downloading %d/%d %v ...\n", chapter.Index, len(chapterList), chapterUrl.String())

		// make file name
		file := filepath.Join(tempdir, fmt.Sprintf("%04d %v.txt", chapter.Index, chapter.Name))
		if fileExist(file) {
			continue
		}

		// download
		html := downloadURL(chapterUrl.String(), config.WebsiteCharset)

		// html to text
		html = getNovelChapterContent(html, config.Content)
		if len(html) < config.Content.IgnoreShortText {
			fmt.Printf("Ignored %s\n", chapter.Name)
			continue
		}

		text := "    " + chapter.Name + "\n\n" + htmlAsText(html)

		// write file
		err = ioutil.WriteFile(file, []byte(text), 0644)
		if err != nil {
			panic(err)
		}
	}

	// zip novel file
	file := filepath.Join(dir, title+".zip")
	fmt.Println()
	fmt.Printf("Archive Zip: %v ...\n", file)
	zipToFile(file, tempdir, config.ZipFilenameCharset)

	// done
	fmt.Println()
	fmt.Println("Completed!")
}
