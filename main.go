package main

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/subchen/go-cli"
)

var (
	BuildVersion   string
	BuildGitRev    string
	BuildGitCommit string
	BuildDate      string
)

var (
	nThreads      int
	nShortChapter int
	directory     string
)

func main() {
	app := cli.NewApp()
	app.Name = "mknovel"
	app.Usage = "Download a novel from URL, transform HTML to TEXT, pack it"
	app.UsageText = "[options] URL"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Flags = []*cli.Flag{
		{
			Name:        "threads",
			Usage:       "parallel threads",
			Placeholder: "num",
			DefValue:    "100",
			Value:       &nThreads,
		},
		{
			Name:        "short-chapter",
			Usage:       "ignore chapter if size is short",
			Placeholder: "size",
			DefValue:    "3000",
			Value:       &nShortChapter,
		}, {
			Name:        "d, directory",
			Usage:       "output directory",
			Placeholder: "dir",
			DefValue:    ".",
			Value:       &directory,
		},
	}

	if BuildVersion != "" {
		app.Version = BuildVersion + "-" + BuildGitRev
		app.BuildGitCommit = BuildGitCommit
		app.BuildDate = BuildDate
	}

	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			c.ShowHelp()
			return
		}

		rawUrl := c.Args()[0]

		bookUrl, err := url.Parse(rawUrl)
		if err != nil || bookUrl.Host == "" {
			c.ShowError(fmt.Errorf("The book url is invalid: %s", bookUrl))
		}

		directory, _ = filepath.Abs(directory)
		downloadNovel(bookUrl, directory, nThreads, nShortChapter)
	}

	app.Run(os.Args)
}
