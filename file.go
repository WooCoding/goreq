package goreq

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// Creates a new file upload http request with optional extra params
func UploadFile(files *Files) (*bytes.Buffer, string, error) {
	var contentType string
	f, err := os.Open(files.path)
	if err != nil {
		return nil, contentType, err
	}
	defer f.Close()

	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	fw, err := w.CreateFormFile(files.name, filepath.Base(files.path))
	if err != nil {
		return nil, contentType, err
	}

	if _, err = io.Copy(fw, f); err != nil {
		return nil, contentType, err
	}

	for key, val := range files.params {
		_ = w.WriteField(key, val)
	}
	contentType = w.FormDataContentType()
	err = w.Close()
	if err != nil {
		return nil, contentType, err
	}
	return body, contentType, nil
}
