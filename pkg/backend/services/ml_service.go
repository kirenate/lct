package services

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"main.go/utils/settings_utils"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"
)

func (r *Service) ProcessWithML(file *multipart.FileHeader) (*[]byte, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	contents, err := file.Open()
	defer contents.Close()

	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	part, err := writer.CreateFormFile("image", filepath.Base(file.Filename))
	if err != nil {
		return nil, errors.Wrap(err, "failed to write field")
	}
	io.Copy(part, contents)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, settings_utils.Settings.MlUrl, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{Timeout: time.Minute * 2}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}
	fmt.Println(resp.StatusCode)
	buf, err := io.ReadAll(resp.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, errors.Wrap(err, "failed to read request body")
	}
	return &buf, nil
}
