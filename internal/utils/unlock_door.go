package utils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func UnlockDoor(doorId string) (bool, error) {
	interval := os.Getenv("ACCESS_CONTROL_CONFIG_INTERVAL")
	apiKey := os.Getenv("ACCESS_CONTROL_CONFIG_API_KEY")
	doorAPIBaseURL := os.Getenv("ACCESS_CONTROL_CONFIG_BASE_URL")
	completeURL, err := url.JoinPath(doorAPIBaseURL, "/api/door/remoteOpenById?doorId=", doorId,
		"&interval=", interval, "&access_token=", apiKey)
	if err != nil {
		return false, err
	}
	httpReq, err := http.NewRequest("POST", completeURL, nil)
	if err != nil {
		return false, err
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}
	resp, err := internalclient.Do(httpReq)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("unexpected response from door endpoint: %d - %s", resp.StatusCode, string(bodyBytes))
	}
	return true, nil
}
