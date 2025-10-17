package services

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io"
	"main.go/utils/settings_utils"
	"mime/multipart"
	"net/http"
	"time"
)

func (r *Service) ProcessWithML(ctx context.Context, contents []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("image", "image.png")
	if err != nil {
		return "", errors.Wrap(err, "failed to write field")
	}

	readerBuf := bytes.NewBuffer(contents)
	_, err = io.Copy(fw, readerBuf)
	if err != nil {
		return "", errors.Wrap(err, "failed to copy contents")
	}

	err = writer.Close()
	if err != nil {
		return "", errors.Wrap(err, "writer close")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, settings_utils.Settings.MlUrl, body)
	if err != nil {
		return "", errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{Timeout: time.Minute * 2}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "request failed")
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", errors.Wrap(err, "failed to read request body")
	}

	zerolog.Ctx(ctx).
		Info().
		Int("resp_len", len(buf)).
		Int("status_code", resp.StatusCode).
		Msg("ml.request.processed")

	return string(buf), nil
}
