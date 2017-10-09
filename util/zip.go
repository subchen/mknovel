package util

import (
	"archive/zip"
	"os"
	"path/filepath"

	"github.com/ungerik/go-dry"
)

func ZipToFile(zipfile string, path string, filenameCharset string) {
	fw, err := os.Create(zipfile)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	w := zip.NewWriter(fw)
	defer w.Close()

	writeDirToZip(w, path, "", filenameCharset)
}

func writeDirToZip(w *zip.Writer, path string, root string, filenameCharset string) {
	dir, err := os.Open(path)
	if err != nil {
		panic(nil)
	}
	defer dir.Close()

	filelist, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	for _, fi := range filelist {
		entryName := fi.Name()
		if root != "" {
			entryName = root + "/" + entryName
		}

		if fi.IsDir() {
			writeDirToZip(w, filepath.Join(path, fi.Name()), entryName, filenameCharset)
			continue
		}

		data, err := dry.FileGetBytes(filepath.Join(path, fi.Name()))
		if err != nil {
			panic(err)
		}

		entryName = string(EncodeBytes([]byte(entryName), filenameCharset))
		f, err := w.Create(entryName)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
