package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"strings"
)

func zipToFile(zipfile string, path string, filenameCharset string) {
	fw, err := os.Create(zipfile)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	w := zip.NewWriter(fw)
	defer w.Close()

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
		if fi.IsDir() {
			continue
		}

		fr, err := os.Open(dir.Name() + "/" + fi.Name())
		if err != nil {
			panic(err)
		}
		defer fr.Close()

		data, err := ioutil.ReadAll(fr)
		if err != nil {
			panic(err)
		}

		name := fi.Name()
		if filenameCharset != "" && strings.ToUpper(charset) != "UTF-8" {
			name = encodeString([]byte(name), filenameCharset)
		}

		f, err := w.Create(name)
		if err != nil {
			panic(err)
		}

		_, err = f.Write(data)
		if err != nil {
			panic(err)
		}
	}
}
