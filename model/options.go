package model

type (
	NovelOptions struct {
		// input
		NovelURL      string
		InputEncoding string

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
		OutputEncoding      string
		ZipFilenameEncoding string

		// debug
		Debug bool
	}
)
