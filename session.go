package goreq

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"strings"
	"time"
	"io"
)

type File struct {
	name   string
	path   string
	params map[string]string
}

type Option struct {
	Header  map[string]string
	Timeout time.Duration
	Proxy   string
	Data    url.Values
	Param   url.Values
	Json    map[string]interface{}
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
	// disable keep alives by default, see this issue https://github.com/parnurzeal/gorequest/issues/75
	s.Transport.DisableKeepAlives = true
	return s
}

//build a new Request
func (s *Session) Request(method, URL string, opt *Option) (*http.Response, error) {

	var (
		body        interface{}
		contentType string
		req         *http.Request
		res         *http.Response
		err         error
	)
	//清空代理
	//s.Transport.Proxy = nil
	// 解析参数
	t := reflect.TypeOf(*opt)
	v := reflect.ValueOf(*opt)
	for i := 0; i < t.NumField(); i++ {
		//fmt.Printf("%s -- %v \n", t.Field(i).Name, v.Field(i).Interface())
		name := t.Field(i).Name
		val := v.Field(i).Interface()
		switch name {
		case "Proxy":
			proxyURL := val.(string)
			if proxyURL == "" {
				s.Transport.Proxy = nil
			} else {
				parsedProxyURL, err := url.Parse(proxyURL)
				if err != nil {
					return nil, err
				}
				s.Transport.Proxy = http.ProxyURL(parsedProxyURL)
			}
		case "Data":
			data := val.(url.Values)
			if len(data) != 0 {
				body = strings.NewReader(data.Encode())
				contentType = "application/x-www-form-urlencoded"
			}
		case "Param":
			param := val.(url.Values)
			if len(param) != 0 {
				var buf bytes.Buffer
				buf.WriteString(URL)
				buf.WriteByte('?')
				buf.WriteString(param.Encode())
				URL = buf.String()
			}
		case "Json":
			jsonData := val.(map[string]interface{})
			if len(jsonData) != 0 {
				bytesData, err := json.Marshal(jsonData)
				if err != nil {
					return nil, err
				}
				body = bytes.NewReader(bytesData)
				contentType = "application/json;charset=UTF-8"
			}
		}
	}

	// 构造请求

	if body != nil {
		req, err = http.NewRequest(method, URL, body.(io.Reader))
	} else {
		req, err = http.NewRequest(method, URL, nil)
	}
	//设置请求头
	for key, val := range opt.Header {
		req.Header.Set(key, val)
	}

	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}

	s.Client.Transport = s.Transport
	s.Client.Timeout = opt.Timeout
	res, err = s.Client.Do(req)
	return res, err
}

//wrap a GET method
func (s Session) Get(URL string, opt *Option) (*http.Response, error) {
	return s.Request("GET", URL, opt)
}

//wrap a POST method
func (s Session) Post(URL string, opt *Option) (*http.Response, error) {
	return s.Request("POST", URL, opt)
}
