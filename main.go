package main

import (
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"

	"github.com/subchen/gstack/cli"
)

const VERSION = "1.0.0"

var (
	BuildVersion   string
	BuildGitCommit string
	BuildDate      string
)

func main() {
	app := cli.NewApp("mknovel", "Download a novel from URL, transform HTML to TEXT, zipped it.")
	app.Flag("--threads", "parallel threads").Default("100")
	app.Flag("-d, --directory", "output directory").Default(".")

	if BuildVersion == "" {
		app.Version = VERSION
	} else {
		app.Version = func() {
			fmt.Printf("Version: %s-%s\n", VERSION, BuildVersion)
			fmt.Printf("Go version: %s\n", runtime.Version())
			fmt.Printf("Git commit: %s\n", BuildGitCommit)
			fmt.Printf("Built: %s\n", BuildDate)
			fmt.Printf("OS/Arch: %s-%s\n", runtime.GOOS, runtime.GOARCH)
		}
	}

	app.Usage = func() {
		fmt.Println("Usage: mknovel [--threads=100] [-d dir] URL")
		fmt.Println("   or: mknovel [ --version | --help ]")
	}

	app.AllowArgumentCount(1, 1)

	app.Execute = func(ctx *cli.Context) {
		nThreads := ctx.Int("--threads")
		dir := ctx.String("-d")
		rawUrl := ctx.Arg(0)

		bookUrl, err := url.Parse(rawUrl)
		if err != nil || bookUrl.Host == "" {
			cli.Fatalf("The book url is invalid: %s", bookUrl)
		}

		dir, _ = filepath.Abs(dir)
		downloadNovel(bookUrl, dir, nThreads)
	}

	app.Run()
}
