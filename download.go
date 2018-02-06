package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadImage(qUrl string) (string, error) {
	var e Err
	var dn DirNav
	var v Validate

	response, err := http.Get(qUrl)
	if err != nil {
		e.setInternalServerError(err)
		return "", err
	}

	defer response.Body.Close()

	addImgDirs()

	name := filepath.Base(qUrl)
	dn.setImgOrigName(name)

	//open a file for writing
	file, err := os.Create(dn.getOrigImgPath())
	if err != nil {
		e.setInternalServerError(err)
		return "", err
	}

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		e.setInternalServerError(err)
		return "", err
	}

	file.Close()

	v.Image(dn.getOrigImgPath())

	return name, err
}
