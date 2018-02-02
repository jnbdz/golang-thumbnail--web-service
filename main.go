package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

type ResponseError struct {
	HttpStatusCode int
	Message        string
}

var isErrorResponse = false
var httpStatusCode = 200
var responseErrorMsgs = []ResponseError{}

func unsetError() {
	isErrorResponse = false
	responseErrorMsgs = []ResponseError{}
	httpStatusCode = 200
}

func setError(message string, httpStatusCode int) {
	isErrorResponse = true
	m := ResponseError{
		httpStatusCode,
		message,
	}
	responseErrorMsgs = append(responseErrorMsgs, m)
}

func setInternalServerError(err error) {
	setError("Internal server error.", 500)
	log.Fatal(err)
}

func getError() []ResponseError {
	return responseErrorMsgs
}

func setHTTPStatusCode() {
	for _, msg := range responseErrorMsgs {
		if httpStatusCode < msg.HttpStatusCode {
			httpStatusCode = msg.HttpStatusCode
		}
	}
}

func getHTTPStatusCode() int {
	return httpStatusCode
}

func sendError(w http.ResponseWriter) {
	b, err := json.Marshal(getError())
	if err != nil {
		setInternalServerError(err)
	}
	if isErrorResponse {
		setHTTPStatusCode()
		http.Error(w, string(b), getHTTPStatusCode())
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

func validateImage(imagePath string) {
	tmpFile, err := os.Open(imagePath)
	if err != nil {
		setInternalServerError(err)
	}

	buff := make([]byte, 512)
	_, err = tmpFile.Read(buff)
	if err != nil {
		setInternalServerError(err)
	}

	contentType := http.DetectContentType(buff)

	tmpFile.Close()

	if !isImage(contentType) {
		// 415 HTTP Status Code: Unsupported Media Type
		setError("The URL is not an image.", 415)
	}
}

func downloadImage(qUrl string) (string, error) {
	response, err := http.Get(qUrl)
	//http.DetectContentType(response)
	if err != nil {
		setInternalServerError(err)
		return "", err
	}

	defer response.Body.Close()

	currPath, err := filepath.Abs("./")
	if err != nil {
		setInternalServerError(err)
	}
	currPath += "/"

	os.MkdirAll(currPath+"imgs/orig", 0777)
	os.MkdirAll(currPath+"imgs/resized", 0777)

	name := filepath.Base(qUrl)

	//open a file for writing
	file, err := os.Create("/tmp/" + name)
	if err != nil {
		setInternalServerError(err)
		return "", err
	}

	// Use io.Copy to just dump the response body to the file. This supports huge files
	_, err = io.Copy(file, response.Body)
	if err != nil {
		setInternalServerError(err)
		return "", err
	}

	file.Close()

	validateImage("/tmp/" + name)

	return name, err
}

func imageResizeName(imageFileName string, width, height int) string {
	infile := imageFileName
	ext := filepath.Ext(infile) // e.g., ".jpg", ".JPEG"
	size := fmt.Sprintf(".%dx%d", width, height)
	return strings.TrimSuffix(infile, ext) + size + ".thumb" + ext
}

func resizeImage(width, height int, infile string) {
	outfile := imageResizeName(infile, width, height)

	file, err := os.Open("/tmp/" + infile)
	if err != nil {
		setInternalServerError(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		setInternalServerError(err)
	}

	file.Close()

	// Documentaion on interpolation opions:
	// https://www.cambridgeincolour.com/tutorials/image-resize-for-web.htm
	m := resize.Thumbnail(uint(width), uint(height), img, resize.Bicubic)

	mb := m.Bounds()
	origWidth := mb.Dx()
	origHeight := mb.Dy()
	paddingTop := (height / 2) - (origHeight / 2)
	paddingLeft := (width / 2) - (origWidth / 2)

	myimageOffset := image.Pt(0, 0)
	myimageRect := image.Rect(0, 0, width, height)
	myimage := image.NewRGBA(myimageRect)
	black := color.RGBA{0, 0, 0, 0}
	draw.Draw(myimage, myimage.Bounds().Add(myimageOffset), &image.Uniform{black}, image.ZP, draw.Src)

	b := myimage.Bounds().Add(image.Pt(-paddingLeft, -paddingTop))
	image3 := image.NewRGBA(b)
	draw.Draw(image3, m.Bounds(), m, image.ZP, draw.Over)

	out, err := os.Create(outfile)
	if err != nil {
		setInternalServerError(err)
	}

	defer out.Close()

	// write new image to file
	jpeg.Encode(out, image3, &jpeg.Options{jpeg.DefaultQuality})
}

func validateUrl(qUrl string) (string, error) {
	_, err := url.ParseRequestURI(qUrl)
	if err != nil {
		setError("Not a url.", 400)
	}

	return qUrl, err
}

func validateSize(sizeStr string, name string) (int, error) {
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		setError(name+" param is not a number.", 400)
	} else if size < 1 {
		setError(name+" param number is too small.", 400)
	}

	return size, err
}

func createThumbnail(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()

	url, _ := validateUrl(urlQuery.Get("url"))
	width, _ := validateSize(urlQuery.Get("width"), "Width")
	height, _ := validateSize(urlQuery.Get("height"), "Height")

	if !isErrorResponse {
		name, _ := downloadImage(url)
		resizeImage(width, height, name)
	}
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		unsetError()

		if r.URL.Path == "/thumbnail" {
			createThumbnail(w, r)
		} else {
			setError("Not Found", 404)
		}
		sendError(w)
	})

	http.ListenAndServe(":3000", nil)
}
