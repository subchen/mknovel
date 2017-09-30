package main

import (
	"io/ioutil"
	"path/filepath"

	"fmt"
	"github.com/go-yaml/yaml"
)

type (
	NovelConfig struct {
		WebsiteCharset string                     `yaml:"website-charset"`
		Title          *NovelTitleConfig          `yaml:"title"`
		Author         *NovelAuthorConfig         `yaml:"author"`
		CoverImage     *NovelCoverImageConfig     `yaml:"cover-image"`
		ChapterIndex   *NovelChapterIndexConfig   `yaml:"chapter-index"`
		ChapterContent *NovelChapterContentConfig `yaml:"chapter-content"`
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

	NovelCoverImageConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
	}

	NovelChapterIndexConfig struct {
		Begin     string `yaml:"begin"`
		End       string `yaml:"end"`
		Regexp    string `yaml:"regexp"`
		NameIndex int    `yaml:"name-index"`
		LinkIndex int    `yaml:"link-index"`
	}

	NovelChapterContentConfig struct {
		Begin string `yaml:"begin"`
		End   string `yaml:"end"`
	}
)

func lookupFile(file string) string {
	filelist := []string{
		filepath.Join(getExecutorDirectory(), file),
		filepath.Join(getExecutorDirectory(), "config", file),
		filepath.Join(getCurrentDirectory(), file),
		filepath.Join(getCurrentDirectory(), "config", file),
		filepath.Join("/etc/mknovel", file),
	}

	for _, f := range filelist {
		if fileExist(f) {
			return f
		}
	}

	panic("file not found: " + file)
}

func loadConfigFile(file string) *NovelConfig {
	configFile := lookupFile(file)
	fmt.Printf("Config File: %s\n", configFile)

	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	config := NovelConfig{}
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		panic(err)
	}

	return &config
}
