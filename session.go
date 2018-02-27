package goreq

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
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
	Client    *http.Client
	Transport *http.Transport
}

//initiate a session
func NewSession() *Session {
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	s := &Session{
		Client:    &http.Client{Jar: jar},
		Transport: &http.Transport{},
	}
	return s
}

//build a new Request
func (s *Session) Request(method, reqURL string, args ...interface{}) (res *http.Response, err error) {

	var body interface{}
	var headers Headers
	var contentType string
	var req *http.Request
	// 解析参数
	for _, arg := range args {
		switch t := arg.(type) {
		case *Headers:
			// 请求头
			headers = *t
		case time.Duration:
			// 超时
			timeout := t
			s.Transport.Dial = func(network, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(network, addr, timeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(timeout))
				return conn, nil
			}
		case Proxy:
			proxyURL := string(t)
			if proxyURL == "" {
				s.Transport.Proxy = nil
			} else {
				parsedProxyURL, err := url.Parse(proxyURL)
				if err != nil {
					return nil, err
				}
				s.Transport.Proxy = http.ProxyURL(parsedProxyURL)
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
	// 构造请求
	if body != nil {
		req, err = http.NewRequest(method, reqURL, body.(io.Reader))
	} else {
		req, err = http.NewRequest(method, reqURL, nil)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	s.Client.Transport = s.Transport
	return s.Client.Do(req)
}

//wrap a GET method
func (s Session) Get(url string, args ...interface{}) (*http.Response, error) {
	return s.Request("GET", url, args...)
}

//wrap a POST method
func (s Session) Post(url string, args ...interface{}) (*http.Response, error) {
	return s.Request("POST", url, args...)
}
