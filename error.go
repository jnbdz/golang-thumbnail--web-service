package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	HttpStatusCode int
	Message        string
}

type Err struct {
	isErrorResponse   bool
	httpStatusCode    int
	responseErrorMsgs []ResponseError
}

func (e *Err) setErrorVars() {
	e.isErrorResponse = false
	e.responseErrorMsgs = []ResponseError{}
	e.httpStatusCode = 200
}

func (e *Err) setError(message string, httpStatusCode int) {
	e.isErrorResponse = true
	m := ResponseError{
		httpStatusCode,
		message,
	}
	e.responseErrorMsgs = append(e.responseErrorMsgs, m)
}

func (e *Err) setInternalServerError(err error) {
	e.setError("Internal server error.", 500)
	log.Fatal(err)
}

func (e *Err) getError() []ResponseError {
	return e.responseErrorMsgs
}

func (e *Err) setHTTPStatusCode() {
	for _, msg := range e.responseErrorMsgs {
		if e.httpStatusCode < msg.HttpStatusCode {
			e.httpStatusCode = msg.HttpStatusCode
		}
	}
}

func (e *Err) getHTTPStatusCode() int {
	return e.httpStatusCode
}

func (e *Err) sendError(w http.ResponseWriter) {
	b, err := json.Marshal(e.getError())
	if err != nil {
		e.setInternalServerError(err)
	}
	if e.isErrorResponse {
		e.setHTTPStatusCode()
		http.Error(w, string(b), e.getHTTPStatusCode())
	}
}
