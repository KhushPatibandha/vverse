package tests

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/KhushPatibandha/vverse/consts"
	db "github.com/KhushPatibandha/vverse/internal/DB"
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

func addMockData() {
	query := `INSERT INTO videos (name, size, duration, created_at) VALUES 
		('e851d391-5ded-40e3-b323-2c44392421d4_1740909320', 1.68618011474609, 10.333333, ?),
		('ff9212f1-981a-479f-99d0-97bc78e83dca_1740909329', 1.1683874130249, 9.373333, ?);`
	_, err := db.ExecCmd(query, "2021-08-01 00:00:00", "2021-08-01 00:00:00")
	if err != nil {
		panic("Failed to add mock data: " + err.Error())
	}

	_ = os.MkdirAll("./uploads", os.ModePerm)
	src1 := "../test_videos/10_work.mp4"
	src2 := "../test_videos/9_work.mp4"
	videoFile1 := "./uploads/e851d391-5ded-40e3-b323-2c44392421d4_1740909320"
	videoFile2 := "./uploads/ff9212f1-981a-479f-99d0-97bc78e83dca_1740909329"

	err = copyFile(src1, videoFile1)
	if err != nil {
		panic("Failed to copy file: " + err.Error())
	}
	err = copyFile(src2, videoFile2)
	if err != nil {
		panic("Failed to copy file: " + err.Error())
	}
}

func cleanMockData() {
	os.Remove("./uploads/e851d391-5ded-40e3-b323-2c44392421d4_1740909320")
	os.Remove("./uploads/ff9212f1-981a-479f-99d0-97bc78e83dca_1740909329")
}
