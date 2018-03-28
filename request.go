package goreq

import "net/http"

func Get(URL string, opt *Option) (*http.Response, error) {
	s := NewSession()
	return s.Get(URL, opt)
}

func Post(URL string, opt *Option) (*http.Response, error) {
	s := NewSession()
	return s.Post(URL, opt)
}
