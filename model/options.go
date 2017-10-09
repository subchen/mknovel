package model

type (
	NovelOptions struct {
		NovelURL string

		// manunal set
		NovelName          string
		NovelAuthor        string
		NovelCoverImageURL string

		// for download
		Threads          int
		ShortChapterSize int

		// output
		OutputFormat        string
		OutputDirectory     string
		TxtEncoding         string
		ZipFilenameEncoding string
	}
)
