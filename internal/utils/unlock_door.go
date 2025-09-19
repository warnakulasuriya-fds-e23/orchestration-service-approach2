package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/valyala/fasthttp"
)

// used fasthttp instead of net/http because it wont send the header as api-key but send it as Api-Key which is not accepted by choreo
func UnlockDoor(doorId string) (bool, error) {
	interval := os.Getenv("ACCESS_CONTROL_CONFIG_INTERVAL")
	apiKey := os.Getenv("ACCESS_CONTROL_CONFIG_API_KEY")
	doorAPIBaseURL := os.Getenv("ACCESS_CONTROL_CONFIG_BASE_URL")

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	uri := fmt.Sprintf("%s/api/door/remoteOpenById?doorId=%s&interval=%s&access_token=%s",
		doorAPIBaseURL, doorId, interval, apiKey)
	req.SetRequestURI(uri)
	req.Header.SetMethod("POST")
	req.Header.Set("api-key", apiKey) // lowercase header

	client := &fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	log.Println("UnlockDoor calling URL:", uri)
	log.Println("Request Headers:", req.Header.String())

	if err := client.Do(req, resp); err != nil {
		return false, err
	}
	if resp.StatusCode() != 200 {
		log.Println("Response from door endpoint:", string(resp.Body()))
		return false, fmt.Errorf("unexpected response from door endpoint: %d - %s", resp.StatusCode(), string(resp.Body()))
	}
	return true, nil
}
