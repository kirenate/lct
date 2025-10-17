package services

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/teadove/teasutils/utils/test_utils"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestProcessWithMl2(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile("/Users/teadove/Downloads/00000010.jpg")
	assert.NoError(t, err)

	r := Service{}
	resp, err := r.ProcessWithML(test_utils.GetLoggedContext(), content)
	require.NoError(t, err)
	test_utils.Pprint(resp)
}

func TestProcessWithMl(t *testing.T) {
	file, err := os.Open("/Users/teadove/Downloads/00000010.jpg")
	assert.NoError(t, err)

	contents, err := io.ReadAll(file)
	assert.NoError(t, err)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	fw, err := writer.CreateFormFile("image", filepath.Base(file.Name()))
	assert.NoError(t, err)

	readerBuf := bytes.NewBuffer(contents)
	_, err = io.Copy(fw, readerBuf)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://158.160.198.84:8081/process", body)
	assert.NoError(t, err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := http.Client{Timeout: time.Minute * 2}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	fmt.Println(resp.StatusCode)
	buf, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	res := string(buf)
	fmt.Println(res)

}
