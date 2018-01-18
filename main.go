package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ResponseError struct {
	Message string
}

func responseErrorMsg(w http.ResponseWriter, message string, httpStatusCode int) {
	m := ResponseError{message}
	b, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	s := string(b)
	fmt.Fprintf(w, s)
}

func isUrl(qUrl string) bool {
	_, err := url.ParseRequestURI(qUrl)
	if err == nil {
		return true
	}
	return false
}

func isNumber(value string) bool {
	if _, err := strconv.Atoi(value); err == nil {
		return true
	}
	return false
}

func validateUrl(qUrl string, w http.ResponseWriter) {
	if !isUrl(qUrl) {
		responseErrorMsg(w, "Not a url.", 500)
	}
}

func validateWidth(width string, w http.ResponseWriter) {
	if !isNumber(width) {
		responseErrorMsg(w, "Width is not a number.", 500)
	}
}

func validateHeight(height string, w http.ResponseWriter) {
	if !isNumber(height) {
		responseErrorMsg(w, "Height is not a number.", 500)
	}
}

func isImage(contentType string) bool {
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

func validateImage(imagePath string, w http.ResponseWriter) {
	tmpFile, err := os.Open(imagePath)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	buff := make([]byte, 512)
	_, err = tmpFile.Read(buff)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	contentType := http.DetectContentType(buff)

	tmpFile.Close()

	if !isImage(contentType) {
		responseErrorMsg(w, "The URL is not an image.", 500)
	}
}

func downloadImage(qUrl string, w http.ResponseWriter) {
	response, err := http.Get(qUrl)
	//http.DetectContentType(response)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	//open a file for writing
	file, err := os.Create("/tmp/asdf.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	validateImage("/tmp/asdf.jpg", w)

	fmt.Fprintf(w, "Success!")
}

func resizeImage(imagePath string, width int, height int) {
	infile := "asdf.jpg"
	ext := filepath.Ext(infile) // e.g., ".jpg", ".JPEG"
	outfile := strings.TrimSuffix(infile, ext) + ".thumb" + ext

	file, err := os.Open("/tmp/" + infile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	out, err := os.Create(outfile)
	if err != nil {
		log.Fatal(err)
	}

	src, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	jpeg.Encode(out, dst, nil)

}

func main() {
	http.HandleFunc("/thumbnail", func(w http.ResponseWriter, r *http.Request) {
		urlQuery := r.URL.Query()

		qUrl := urlQuery.Get("url")
		qWidth := urlQuery.Get("width")
		qHeight := urlQuery.Get("height")

		validateUrl(qUrl, w)
		validateWidth(qWidth, w)
		validateHeight(qHeight, w)

		//fmt.Fprintf(w, "Hello Gopher thumbnail\n")
		fmt.Fprintf(w, qUrl+"\n")
		fmt.Fprintf(w, qWidth+"\n")
		fmt.Fprintf(w, qHeight+"\n")

		width, err := strconv.Atoi(qWidth)
		if err != nil {
			log.Fatal(err)
		}

		height, err := strconv.Atoi(qHeight)
		if err != nil {
			log.Fatal(err)
		}

		resizeImage(width, height)

		//downloadImage(qUrl, w)

		w.Header().Add("Server", "Thumbnail Web Service Server")
	})

	http.ListenAndServe(":3000", nil)
}
