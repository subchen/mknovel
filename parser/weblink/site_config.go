package weblink

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/subchen/go-stack"
	"github.com/subchen/mknovel/model"
	"github.com/ungerik/go-dry"
)

type (
	SiteConfig struct {
		WebsiteCharset string                    `yaml:"website-charset"`
		Title          *SiteTitleConfig          `yaml:"title"`
		Author         *SiteAuthorConfig         `yaml:"author"`
		CoverImage     *SiteCoverImageConfig     `yaml:"cover-image"`
		ChapterIndex   *SiteChapterIndexConfig   `yaml:"chapter-index"`
		ChapterContent *SiteChapterContentConfig `yaml:"chapter-content"`
	}

	SiteTitleConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
	}

	SiteAuthorConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
	}

	SiteCoverImageConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
	}

	SiteChapterIndexConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
		LinkIndex int    `yaml:"link-index"`
	}

	SiteChapterContentConfig struct {
		Begin string `yaml:"begin"`
		End   string `yaml:"end"`
	}
)

func loadConfig(file string) *SiteConfig {
	dirs := []string{
		gstack.ProcessGetPWD(),
		filepath.Join(gstack.ProcessGetPWD(), "site-config"),
		gstack.ProcessGetBinDir(),
		filepath.Join(gstack.ProcessGetBinDir(), "site-config"),
		"/etc/mknovel",
	}

	configFile, found := dry.FileFind(dirs, file)
	if !found {
		panic("file not found: " + file)
	}

	fmt.Printf("Site Config File: %s\n", configFile)
	bytes, err := ioutil.ReadFile(configFile)
	dry.PanicIfErr(err)

	config := SiteConfig{}
	err = yaml.Unmarshal(bytes, &config)
	dry.PanicIfErr(err)

	return &config
}

func getNovelTitle(html string, config *SiteTitleConfig) string {
	if config.Begin != "" && config.End != "" {
		html = gstack.StringBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindStringSubmatch(html)
	title := matches[config.NameIndex]
	return strings.TrimSpace(title)
}

func getNovelAuthor(html string, config *SiteAuthorConfig) string {
	if config.Begin != "" && config.End != "" {
		html = gstack.StringBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindStringSubmatch(html)
	title := matches[config.NameIndex]
	return strings.TrimSpace(title)
}

func getNovelCoverImage(html string, config *SiteCoverImageConfig) string {
	if config.Begin != "" && config.End != "" {
		html = gstack.StringBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindStringSubmatch(html)
	title := matches[config.NameIndex]
	return strings.TrimSpace(title)
}

func getNovelChapterList(html string, config *SiteChapterIndexConfig) []*model.NovelChapter {
	if config.Begin != "" && config.End != "" {
		html = gstack.StringBetween(html, config.Begin, config.End)
	}

	re := regexp.MustCompile(config.Regexp)
	matches := re.FindAllStringSubmatch(html, -1)
	chapterList := make([]*model.NovelChapter, 0, len(matches))

	for i, match := range matches {
		chapter := &model.NovelChapter{
			ID:    fmt.Sprintf("%04d", i+1),
			Index: i + 1, // from base 1
			Name:  strings.TrimSpace(match[config.NameIndex]),
			URL:   match[config.LinkIndex],
		}
		chapterList = append(chapterList, chapter)
	}
	return chapterList
}

func getNovelChapterContent(html string, config *SiteChapterContentConfig) string {
	html = gstack.StringBetween(html, config.Begin, config.End)
	return html
}
