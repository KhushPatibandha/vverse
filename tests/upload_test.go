package tests

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/vverse/consts"
)

func TestUploadVideo(t *testing.T) {
	resetDatabase()
	filePath := "../test_videos/10_work.mp4"
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	fileStat, err := file.Stat()
	assert.NoError(t, err)

	body := make([]byte, fileStat.Size())
	_, err = file.Read(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://localhost:"+consts.LOCALHOSTPORT+"/api/v1/video", bytes.NewReader(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "someCrazySecureToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedResponse := `{"Code":200,"Message":"Video uploaded!! Video Id: 1, use this for further operations"}`
	assert.JSONEq(t, expectedResponse, string(respBody))
}
