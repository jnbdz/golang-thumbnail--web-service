package main

import (
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Validate struct{}

func (v *Validate) isImage(contentType string) bool {
	switch contentType {
	case
		"image/jpg",
		"image/jpeg",
		"image/png",
		"image/gif":
		return true
	}
	return false
}

func (v *Validate) Image(imagePath string) {
	var e Err
	tmpFile, err := os.Open(imagePath)
	if err != nil {
		e.setInternalServerError(err)
	}

	buff := make([]byte, 512)
	_, err = tmpFile.Read(buff)
	if err != nil {
		e.setInternalServerError(err)
	}

	contentType := http.DetectContentType(buff)

	tmpFile.Close()

	if !v.isImage(contentType) {
		// 415 HTTP Status Code: Unsupported Media Type
		e.setError("The URL is not an image.", 415)
	}
}

func (v *Validate) Url(qUrl string) (string, error) {
	var e Err
	_, err := url.ParseRequestURI(qUrl)
	if err != nil {
		e.setError("Not a url.", 400)
	}

	return qUrl, err
}

func (v *Validate) Size(sizeStr string, name string) (int, error) {
	var e Err
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		e.setError(name+" param is not a number.", 400)
	} else if size < 1 {
		e.setError(name+" param number is too small.", 400)
	}

	return size, err
}
