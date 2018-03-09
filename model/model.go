package model

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

type (
	Request struct {
		Method string            `json:"method"`
		URL    string            `json:"url"`
		Proto  string            `json:"version"`
		Header map[string]string `json:"headers"`
		Body   []byte            `json:"body"`
	}
)

func NewFromHTTPRequest(r *http.Request) (*Request, error) {
	method := r.Method
	url := r.URL.String()
	proto := r.Proto
	header := fromHTTPHeader(r.Header)
	body, err := getHTTPBody(r)
	if err != nil {
		return nil, err
	}
	return &Request{Method: method, URL: url, Proto: proto, Header: header, Body: body}, nil
}

func fromHTTPHeader(header http.Header) map[string]string {
	h := map[string]string{}
	for k, v := range header {
		h[k] = strings.Join(v, " ")
	}
	return h
}

func getHTTPBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}
