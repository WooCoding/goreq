package goreq

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"encoding/json"
	"bytes"
)

type Headers map[string]string
type Json map[string]interface{}
type Data url.Values
type Params url.Values
type Proxy string
type Files struct {
	name   string
	path   string
	params map[string]string
}

type Session struct {
	Client *http.Client
}

//initiate a session
func NewSession() *Session {
	s := &Session{}
	s.Client = &http.Client{}
	Jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	s.Client.Jar = Jar
	return s
}

//build a new Request
func (s *Session) Request(method, reqURL string, args ...interface{}) (*http.Request, error) {

	var body interface{}
	var headers Headers
	var contentType string
	var err error

	for _, arg := range args {
		switch t := arg.(type) {
		case *Headers:
			headers = *t
		case Proxy:
			urlProxy, err := url.Parse(string(t))
			if err != nil {
				return nil, err
			}
			s.Client.Transport = &http.Transport{
				Proxy: http.ProxyURL(urlProxy),
			}
		case *Data:
			body = strings.NewReader(url.Values(*t).Encode())
			contentType = "application/x-www-form-urlencoded"
		case *Json:
			bytesData, err := json.Marshal(*t)
			if err != nil {
				return nil, err
			}
			body = bytes.NewReader(bytesData)
			contentType = "application/json;charset=UTF-8"
		case *Params:
			var buf bytes.Buffer
			buf.WriteString(reqURL)
			buf.WriteByte('?')
			buf.WriteString(url.Values(*t).Encode())
			reqURL = buf.String()
		case *Files:
			body, contentType, err = UploadFile(t)
			if err != nil {
				return nil, err
			}
		}
	}

	if body != nil {
		req, err := http.NewRequest(method, reqURL, body.(io.Reader))
		for key, val := range headers {
			req.Header.Set(key, val)
		}
		if contentType != "" {
			req.Header.Add("Content-Type", contentType)
		}
		return req, err
	} else {
		req, err := http.NewRequest(method, reqURL, nil)
		for key, val := range headers {
			req.Header.Set(key, val)
		}
		if contentType != "" {
			req.Header.Add("Content-Type", contentType)
		}
		return req, err
	}
}

//wrap a GET method
func (s Session) Get(url string, args ...interface{}) (*http.Response, error) {
	req, _ := s.Request("GET", url, args...)
	return s.Client.Do(req)
}

//wrap a POST method
func (s Session) Post(url string, args ...interface{}) (*http.Response, error) {
	req, _ := s.Request("POST", url, args...)
	return s.Client.Do(req)
}
