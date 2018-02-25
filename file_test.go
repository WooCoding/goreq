package goreq

import (
	"testing"
	"io/ioutil"
	"github.com/tidwall/gjson"
)

func TestUploadFile(t *testing.T) {
	files := &Files{
		name: "file",
		path: "./file_test.txt",
		params: map[string]string{
			"title":       "My Document",
			"author":      "Matt Aimonetti",
			"description": "A document with all the Go programming language secrets",
		},
	}
	s := NewSession()
	res, err := s.Post("http://httpbin.org/post", files)
	if err != nil {
		t.Error("fail to get a Response", err)
	}
	bodyByte, _ := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	value := gjson.Get(string(bodyByte), "files.file")
	if value.String() != "for test" {
		t.Error("Expected for test, Got ", value.String())
	}
}
