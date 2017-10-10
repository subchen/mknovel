package generator

import (
	"github.com/subchen/mknovel/generator/epub"
	"github.com/subchen/mknovel/generator/txt"
	"github.com/subchen/mknovel/generator/zip"
	"github.com/subchen/mknovel/model"
	"github.com/ungerik/go-dry"
)

var (
	SUPPORTED_FORMATS = []string{
		"txt", "zip", "epub",
	}
)

func ValidateOutputFormat(format string) {
	if !dry.StringInSlice(format, SUPPORTED_FORMATS) {
		panic("Unsupported format: " + format)
	}
}

func PackageNovel(novel *model.Novel, opts *model.NovelOptions) {
	switch opts.OutputFormat {
	case "txt":
		txt.PackageNovelAsTXT(novel, opts.OutputDirectory, opts.OutputEncoding, opts.Debug)
	case "zip":
		zip.PackageNovelAsZIP(novel, opts.OutputDirectory, opts.OutputEncoding, opts.ZipFilenameEncoding, opts.Debug)
	case "epub":
		epub.PackageNovelAsEPUB(novel, opts.OutputDirectory, opts.Debug)
	}
}
