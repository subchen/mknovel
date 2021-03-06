package model

type (
	NovelOptions struct {
		// input
		NovelURL      string
		InputEncoding string

		// txt parse
		AutoChapterGroup bool

		// manunal set
		NovelName          string
		NovelAuthor        string
		NovelCoverImageURL string

		// for download
		Threads          int
		ShortChapterSize int

		// output
		OutputFormat    string
		OutputDirectory string
		OutputEncoding  string

		// debug
		Debug bool
	}
)
