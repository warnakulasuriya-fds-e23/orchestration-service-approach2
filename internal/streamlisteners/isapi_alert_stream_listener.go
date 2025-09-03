package streamlisteners

import (
	"fmt"
	"io"
	"log"
	"mime" // Correct package for ParseMediaType
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/icholy/digest"
	"github.com/tidwall/gjson"
)

// StartAlertStreamListener starts a goroutine to listen to the alert stream
func StartAlertStreamListener(endpoint string) {
	go func() {
		for {
			log.Println("Connecting to alert stream...")
			err := listenForAlerts(endpoint)
			if err != nil {
				log.Printf("Alert stream listener failed with error: %v. Retrying in 5 seconds...", err)
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func listenForAlerts(endpoint string) error {
	// Get credentials from environment variables
	username := os.Getenv("CAM_DIGEST_AUTH_USERNAME")
	password := os.Getenv("CAM_DIGEST_AUTH_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("digest auth username or password not set in environment variables")
	}

	// Create an http.Client with a custom Transport for Digest Auth
	client := &http.Client{
		Transport: &digest.Transport{
			Username: username,
			Password: password,
		},
	}

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from endpoint: %d", resp.StatusCode)
	} else {
		log.Println("Successfully connected to alert stream")
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/") {
		return fmt.Errorf("unexpected Content-Type: %s", contentType)
	}

	// Use mime.ParseMediaType to get the boundary from the Content-Type header
	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	boundary := params["boundary"]

	reader := multipart.NewReader(resp.Body, boundary)

	for {
		part, err := reader.NextPart()
		if err != nil {
			if err == io.EOF {
				return nil // Stream ended gracefully
			}
			return err // Connection error
		}

		switch part.Header.Get("Content-Type") {
		case "application/json":
			// Read the entire part body into a byte slice
			body, err := io.ReadAll(part)
			if err != nil {
				log.Printf("Error reading JSON part body: %v", err)
				continue
			}

			// Use gjson to extract the eventType
			eventType := gjson.GetBytes(body, "eventType").String()

			// Check if this is an event your service needs to handle
			if eventType == "alarmResult" {
				// Use gjson to check the error message in the first array element
				errorMessage := gjson.GetBytes(body, "alarmResult.0.errorMsg").String()

				if errorMessage == "ok" {
					// Process the event as the error message is "ok"
					log.Printf("Received and will process an 'alarmResult' event. Event body: %s", string(body))
					name := gjson.GetBytes(body, "alarmResult.0.faces.0.identify.0.candidate.0.reserve_field.name").String()
					deviceId := gjson.GetBytes(body, "alarmResult.0.targetAttrs.deviceId").String()
					log.Println("The name value of identified employee ", name)
					log.Println("The device ID of identified employee ", deviceId)
					// extracting userid pres
					// nameData := strings.Fields(name)
					// var userId string
					// if len(nameData) == 2 {
					// 	userId = nameData[1]
					// } else if len(nameData) == 1 {
					// 	userId = nameData[0]
					// } else {
					// 	log.Println("Received an 'alarmResult' event with an unexpected name format. Dropping.")
					// } // TODO: Initialize http client at outgoingFingerprintController Startup
					// tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
					// internalclient := &http.Client{Transport: tr}
					// httpReq, err := http.NewRequest("POST", "http://localhost:", userId), nil)

				} else {
					// The error message is not "ok"
					log.Printf("Received an 'alarmResult' event with a non-ok message: '%s'. Dropping.", errorMessage)
				}
			} else {
				// Drop or log other types of events
				log.Printf("Received JSON event of type '%s'. Dropping.", eventType)
			}
		case "image/jpeg":
			log.Println("Received an image part. Discarding...")
			io.Copy(io.Discard, part)
		default:
			log.Printf("Received part with unhandled Content-Type: %s. Discarding...", part.Header.Get("Content-Type"))
			io.Copy(io.Discard, part)
		}
	}
}
