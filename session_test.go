package goreq

import (
	"io/ioutil"
	"testing"

	"github.com/tidwall/gjson"
	"net/url"
	"time"
)

func TestGetParams(t *testing.T) {

	opt := &Option{
		Header:map[string]string{
			"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.167 Safari/537.36",
		},
		Param:url.Values{"arg": {"param"}},
		Timeout:10 * time.Second,
	}
	s := NewSession()
	res, err := s.Get("http://httpbin.org/get", opt)
	if err != nil {
		t.Error("fail to get a Response", err)
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	value := gjson.Get(string(bodyByte), "url")
	if value.String() != "http://httpbin.org/get?arg=param" {
		t.Error("Expected http://httpbin.org/get?arg=param, Got ", value.String())
	}
}

func TestGetHeaders(t *testing.T) {

	opt := &Option{
		Header:map[string]string{
			"User-Agent": "goreq",
		},
	}
	s := NewSession()
	res, err := s.Get("http://httpbin.org/get", opt)
	if err != nil {
		t.Error("fail to get a Response", err)
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	value := gjson.Get(string(bodyByte), "headers.User-Agent")
	if value.String() != "goreq" {
		t.Error("Expected goreq, Got ", value.String())
	}
}

func TestPostData(t *testing.T) {

	opt := &Option{
		Data:url.Values{"key": {"data"}},
	}

	s := NewSession()
	res, err := s.Post("http://httpbin.org/post", opt)
	if err != nil {
		t.Error("fail to get a Response", err)
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	value := gjson.Get(string(bodyByte), "form.key")
	if value.String() != "data" {
		t.Error("Expected data, Got ", value.String())
	}
}

func TestPostJson(t *testing.T) {
	// list
	//json := &Json{
	//	"key":[]string{"json", "list"},
	//}
	json := map[string]interface{}{
		"key": "json",
	}

	opt := &Option{
		Json:json,
	}

	s := NewSession()
	res, err := s.Post("http://httpbin.org/post", opt)
	if err != nil {
		t.Error("fail to get a Response", err)
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	value := gjson.Get(string(bodyByte), "json.key")
	if value.String() != "json" {
		t.Error("Expected json, Got ", value.String())
	}
}

//func TestProxy(t *testing.T) {
//	proxy := Proxy("http://221.7.255.168:80")
//	s := NewSession()
//	res, err := s.Get("http://httpbin.org/get", proxy)
//	if err != nil {
//		t.Error("fail to get a Response", err)
//	}
//	bodyByte, _ := ioutil.ReadAll(res.Body)
//	defer res.Body.Close()
//	value := gjson.Get(string(bodyByte), "origin")
//	if value.String() != "221.7.255.168" {
//		t.Error("Expected 221.7.255.168, Got ", value.String())
//	}
//}
