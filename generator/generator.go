package generator

import (
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
	switch opts.OutputFormat {
	case "epub":
		epub.PackageNovelAsEPUB(novel, opts.OutputDirectory, opts.Debug)
	case "txt":
		txt.PackageNovelAsTXT(novel, opts.OutputDirectory, opts.OutputEncoding, opts.Debug)
	}
}
