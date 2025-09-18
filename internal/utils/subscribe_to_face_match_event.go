package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
)

func SubscribeToFaceMatchEvent() {
	eventDestinationURL, err := url.JoinPath(os.Getenv("ORCHESTRATION_SERVICE_BASE_URL"), "/api/v1/event-receive/receive-face-match-event")
	if err != nil {
		log.Fatal("error while joining event destination URL path:", err)
		return
	}
	requestObj := models.SubscribeToFaceMatchRequest{
		EventTypes: []int{131659},
		EventDest:  eventDestinationURL,
	}
	jsonReqBody, err := json.Marshal(requestObj)
	if err != nil {
		log.Fatal("error while marshalling subscribe to face match request object to JSON:", err)
		return
	}
	subscribeToFaceMatchUrl, err := url.JoinPath("https://", os.Getenv("HCP_IP_ADDRESS"), "/artemis/api/eventService/v1/eventSubscriptionByEventTypes")
	if err != nil {
		log.Fatal("error while joining subscribe to face match URL path:", err)
		return
	}
	req, err := http.NewRequest("POST", subscribeToFaceMatchUrl, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		log.Fatal("error while creating new request to internal authorization endpoint:", err)
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-CA-Key", os.Getenv("HCP_OPENAPI_USER_KEY"))
	requestURLPath := "/artemis/api/resource/v1/person/personId/personInfo"
	generatedSignature, err := GenerateSignatureForHcpOpenapi("POST", req.Header, requestURLPath, os.Getenv("HCP_OPENAPI_USER_SECRET"))
	if err != nil {
		log.Fatal("error while generating signature for HCP OpenAPI:", err)
		return
	}
	req.Header.Set("X-CA-Signature", generatedSignature)
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}
	maxAttempts := 5
	for {

		resp, err := internalclient.Do(req)
		if err != nil {
			maxAttempts--
			if maxAttempts <= 0 {
				log.Fatal("error while sending request to subscribe to face match event endpoint max attempts exceeded:", err)
				return
			}
			log.Println("error while sending request to subscribe to face match event endpoint. Retrying...", err)
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			maxAttempts--
			if maxAttempts <= 0 {
				log.Fatal("unexpected status code from subscribe to face match event endpoint max attempts exceeded:", resp.StatusCode)
				return
			}
			log.Println("unexpected status code from subscribe to face match event endpoint. Retrying...", resp.StatusCode)
			continue
		} else {
			log.Println("Successfully subscribed to FaceMatch events:")
			break
		}
	}

}
