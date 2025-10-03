package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"main.go/schemas"
	"net/http"
)

func (r *Service) ProcessWithML(u string, path string) (*[]schemas.TextJson, error) {
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

	var text *schemas.Text
	var textParsed []schemas.TextJson
	err = json.Unmarshal(buf, &text)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal text")
	}
	for i, v := range text.RecScores {
		textParsed[i].Confidence = v
	}
	for i, v := range text.RecPolys {
		textParsed[i].X1 = v[0]
		textParsed[i].Y1 = v[1]
		textParsed[i].X2 = v[2]
		textParsed[i].Y2 = v[3]
		textParsed[i].X3 = v[4]
		textParsed[i].Y3 = v[5]
		textParsed[i].X4 = v[6]
		textParsed[i].Y4 = v[7]
	}
	for i, v := range text.RecTexts {
		textParsed[i].Text = v
	}
	fmt.Println(textParsed)
	return &textParsed, nil
}
