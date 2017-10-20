package generator

import (
	"os"
	"path/filepath"

	"github.com/subchen/mknovel/generator/epub"
	"github.com/subchen/mknovel/generator/txt"
	"github.com/subchen/mknovel/model"
	"github.com/ungerik/go-dry"
)

var (
	SUPPORTED_FORMATS = []string{
		"txt", "epub",
	}
)

func ValidateOutputFormat(format string) {
	if !dry.StringInSlice(format, SUPPORTED_FORMATS) {
		panic("Unsupported format: " + format)
	}
}

func PackageNovel(novel *model.Novel, opts *model.NovelOptions) {
	outputDirectory, _ := filepath.Abs(opts.OutputDirectory)
	if !dry.FileExists(outputDirectory) {
		err := os.MkdirAll(outputDirectory, 0755)
		dry.PanicIfErr(err)
	}

	switch opts.OutputFormat {
	case "epub":
		epub.PackageNovelAsEPUB(novel, outputDirectory, opts.Debug)
	case "txt":
		txt.PackageNovelAsTXT(novel, outputDirectory, opts.OutputEncoding, opts.Debug)
	}
}
