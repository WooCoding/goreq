package goreq

import "net/http"

func Get(url string, args ...interface{}) (*http.Response, error) {
	s := NewSession()
	return s.Get(url, args...)
}

func Post(url string, args ...interface{}) (*http.Response, error) {
	s := NewSession()
	return s.Post(url, args...)
}
