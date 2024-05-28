package api_test

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wI2L/jsondiff"
)

func TestGetBoards(t *testing.T) {
	resp, err := makeRequest("GET", "/trellode-api/v1/boards", "", "olivier.delobre@gmail.com")
	assert.Equal(t, err, nil)
	responseBytes, _ := io.ReadAll(resp.Body)
	responseContent := string(responseBytes)
	assert.Equal(t, 200, resp.StatusCode)
	// compare with reference response
	equals, diffs, err := compareResponses(responseContent, "GET_list_olivier.delobre.json")
	assert.Nil(t, err)
	if !equals {
		fmt.Printf("JSON is different: %s\n", string(diffs))
	}
	assert.Equal(t, true, equals)
}

func makeRequest(verb string, url string, payload string, userId string) (*http.Response, error) {
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{Transport: customTransport}

	bodyReader := bytes.NewReader([]byte(payload))

	req, _ := http.NewRequest(verb, "http://localhost:8080"+url, bodyReader)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Krakend-UserType", "service")
	req.Header.Add("X-Krakend-UserId", userId)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error calling %s: %s", "http://localhost:8080"+url, err.Error())
		return nil, err
	}

	return resp, nil
}

func compareResponses(actual string, referenceFilename string) (bool, string, error) {
	var v1, v2 interface{}

	// marshal actual response date
	json.Unmarshal([]byte(actual), &v1)

	// read 'expected' data from file
	pwd, _ := os.Getwd()
	rootPath := strings.ReplaceAll(pwd, "internal/api", "")
	b, err := os.ReadFile(rootPath + "assets/tests/" + referenceFilename)
	if err != nil {
		return false, "", err
	}
	json.Unmarshal(b, &v2)

	if reflect.DeepEqual(v1, v2) {
		return true, "", nil
	} else {
		patch, err := jsondiff.CompareJSON([]byte(actual), b)
		if err != nil {
			return false, "", err
		}
		diffs, err := json.MarshalIndent(patch, "", "    ")
		if err != nil {
			return false, "", err
		}
		//fmt.Printf("%s\n", string(diffs))
		return false, string(diffs), nil
	}
}
