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
	nThreads            int
	nShortChapter       int
	trimTrailingAd      bool
	txtEncoding         string
	zipFilenameEncoding string
	outputFormat        string
	outputDirectory     string
)

func main() {
	app := cli.NewApp()
	app.Name = "mknovel"
	app.Usage = "Download a novel from URL and output txt/zip/epub format"
	app.UsageText = "[options] URL"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Flags = []*cli.Flag{
		{
			Name:        "threads",
			Usage:       "parallel threads",
			Placeholder: "num",
			DefValue:    "100",
			Value:       &nThreads,
		}, {
			Name:        "short-chapter",
			Usage:       "ignore chapter if size is short",
			Placeholder: "size",
			DefValue:    "3000",
			Value:       &nShortChapter,
		}, {
			Name:     "trim-trailing-ad",
			Usage:    "remove ad in trailing by keywords",
			DefValue: "false",
			Value:    &trimTrailingAd,
		}, {
			Name:     "txt-encoding",
			Usage:    "encoding for txt file",
			DefValue: "GBK",
			Value:    &txtEncoding,
		}, {
			Name:     "zip-filename-encoding",
			Usage:    "encoding for zip file name",
			DefValue: "GBK",
			Value:    &zipFilenameEncoding,
		}, {
			Name:     "format",
			Usage:    "output file format (txt, zip, epub)",
			DefValue: "epub",
			Value:    &outputFormat,
		}, {
			Name:        "d, directory",
			Usage:       "output directory",
			Placeholder: "dir",
			DefValue:    ".",
			Value:       &outputDirectory,
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

		// dir to abs
		outputDirectory, _ = filepath.Abs(outputDirectory)

		// check book url
		rawUrl := c.Args()[0]
		bookUrl, err := url.Parse(rawUrl)
		if err != nil || bookUrl.Host == "" {
			c.ShowError(fmt.Errorf("Novel URL is invalid: %s", bookUrl))
		}

		// create novel object
		novel := newNovel(bookUrl, outputDirectory)

		downloadNovel(novel, nThreads, trimTrailingAd)

		filterNovelShortChapters(novel, nShortChapter)

		switch outputFormat {
		case "txt":
			packageNovelAsTXT(novel, outputDirectory, txtEncoding)
		case "zip":
			packageNovelAsZIP(novel, outputDirectory, txtEncoding, zipFilenameEncoding)
		case "epub":
			packageNovelAsEPUB(novel, outputDirectory)
		default:
			c.ShowError(fmt.Errorf("Unsupported format: %s", outputFormat))
		}

		// done
		fmt.Println()
		fmt.Println("Completed!")
	}

	app.Run(os.Args)
}
