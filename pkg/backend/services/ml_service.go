package services

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"main.go/schemas"
	"net/http"
)

func (r *Service) ProcessWithML(u string, path string) (*schemas.Text, error) {
	client := http.Client{}
	body := bytes.NewReader([]byte(path))
	req, err := http.NewRequest("POST", u, body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	var attrs *schemas.Text
	err = json.Unmarshal(buf, &attrs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal text")
	}

	return attrs, nil
}
