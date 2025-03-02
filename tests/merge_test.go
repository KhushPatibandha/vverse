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

func TestMergeVideo(t *testing.T) {
	resetDatabase()
	addMockData()

	videoID1 := 1
	videoID2 := 2

	mergeURL := "http://localhost:" + consts.LOCALHOSTPORT + "/api/v1/merge"
	params := url.Values{}
	params.Add("v1", strconv.Itoa(videoID1))
	params.Add("v2", strconv.Itoa(videoID2))

	req, err := http.NewRequest(http.MethodPost, mergeURL+"?"+params.Encode(), nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "someCrazySecureToken")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	expectedResponse := `{"Code":200,"Message":"Videos merged!! New merged video Id: 3, use this for further operations"}`
	assert.JSONEq(t, expectedResponse, string(respBody))
	cleanMockData()
}
