package main

import (
	"net/http"
)

func createThumbnail(w http.ResponseWriter, r *http.Request) {
	var e Err
	var v Validate

	urlQuery := r.URL.Query()

	url, _ := v.Url(urlQuery.Get("url"))
	width, _ := v.Size(urlQuery.Get("width"), "Width")
	height, _ := v.Size(urlQuery.Get("height"), "Height")

	if !e.isErrorResponse {
		name, _ := downloadImage(url)
		resizeImage(width, height, name)
	}
}

func main() {
	var e Err
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.Header.Set("Content-Type", "application/json")
		e.setErrorVars()

		if r.URL.Path == "/thumbnail" {
			createThumbnail(w, r)
		} else {
			e.setError("Not Found", 404)
		}
		e.sendError(w)
	})

	http.ListenAndServe(":3000", nil)
}
