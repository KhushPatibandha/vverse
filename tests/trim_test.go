package tests

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/vverse/consts"
)

func TestTrimVideo(t *testing.T) {
	videoID := 1

	startTime := 2
	endTime := 8

	trimURL := "http://localhost:" + consts.LOCALHOSTPORT + "/api/v1/trim"
	params := url.Values{}
	params.Add("id", strconv.Itoa(videoID))
	params.Add("s", strconv.Itoa(startTime))
	params.Add("e", strconv.Itoa(endTime))

	req, err := http.NewRequest(http.MethodPut, trimURL+"?"+params.Encode(), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "someCrazySecureToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedResponse := `{"Code":200,"Message":"Video trimmed successfully, you can get the trimmed video with the same ID"}`
	assert.JSONEq(t, expectedResponse, string(respBody))
}
