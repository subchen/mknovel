package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/subchen/go-cli"
	"github.com/subchen/mknovel/generator"
	"github.com/subchen/mknovel/model"
	"github.com/subchen/mknovel/parser/weblink"
)

var (
	BuildVersion   string
	BuildGitRev    string
	BuildGitCommit string
	BuildDate      string
)

func main() {
	opts := new(model.NovelOptions)

	app := cli.NewApp()
	app.Name = "mknovel"
	app.Usage = "Download a novel from URL and output txt/zip/epub format"
	app.UsageText = "[options] URL"
	app.Authors = "Guoqiang Chen <subchen@gmail.com>"

	app.Flags = []*cli.Flag{
		{
			Name:  "novel-name",
			Usage: "name of novel",
			Value: &opts.NovelName,
		}, {
			Name:  "novel-author",
			Usage: "author of novel",
			Value: &opts.NovelAuthor,
		}, {
			Name:  "novel-cover-image",
			Usage: "cover image file or url",
			Value: &opts.NovelCoverImageURL,
		}, {
			Name:        "threads",
			Usage:       "parallel threads",
			Placeholder: "num",
			DefValue:    "100",
			Value:       &opts.Threads,
		}, {
			Name:        "short-chapter-size",
			Usage:       "skip chapter if size is short",
			Placeholder: "size",
			DefValue:    "3000",
			Value:       &opts.ShortChapterSize,
		}, {
			Name:     "format",
			Usage:    "output file format (txt, zip, epub)",
			DefValue: "epub",
			Value:    &opts.OutputFormat,
		}, {
			Name:        "d, directory",
			Usage:       "output directory",
			Placeholder: "dir",
			DefValue:    ".",
			Value:       &opts.OutputDirectory,
		}, {
			Name:     "txt-encoding",
			Usage:    "encoding for output txt file",
			DefValue: "GBK",
			Value:    &opts.TxtEncoding,
		}, {
			Name:     "zip-filename-encoding",
			Usage:    "encoding for output file name in zip",
			DefValue: "GBK",
			Value:    &opts.ZipFilenameEncoding,
		},
	}

	// set compiler version
	if BuildVersion != "" {
		app.Version = BuildVersion + "-" + BuildGitRev
	}
	app.BuildGitCommit = BuildGitCommit
	app.BuildDate = BuildDate

	// cli action
	app.Action = func(c *cli.Context) {
		if c.NArg() == 0 {
			c.ShowHelp()
			return
		}

		// validate
		generator.ValidateOutputFormat(opts.OutputFormat)

		// set novel file or url
		opts.NovelURL = c.Args()[0]

		// create novel object
		novel := model.NewNovel(opts)

		// download or parse
		if strings.Contains(opts.NovelURL, "://") {
			// download from url
			weblink.StartDownload(novel, opts.Threads, opts.ShortChapterSize)
		} else {
			// import from a local txt file
			//txtfile.ImportNovel(novel)
		}

		// output novel
		generator.PackageNovel(novel, opts)

		// done
		fmt.Println()
		fmt.Println("Completed!")
	}

	app.Run(os.Args)
}
