package v1_controllers

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/internalkey"
)

type EventReceiveController struct{}

func (erc *EventReceiveController) ReceiveFaceMatchEvent(c *gin.Context) {
	BodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("❌ ERROR: Failed to read request body! Reason: %v", err)
		c.JSON(500, gin.H{"error while reading request body": err.Error()})
		return
	}

	gjsonResult := gjson.ParseBytes(BodyBytes)
	// Check if the event is a FaceMatch event with a recognized human_id
	if gjsonResult.Get("method").String() == "OnEventNotify" &&
		gjsonResult.Get("params.ability").String() == "event_face_match" &&
		gjsonResult.Get("params.events.0.eventType").String() == "131659" &&
		gjsonResult.Get("params.events.0.data.alarmResult.faces.identify.candidate.human_id").String() != "-1" {

		log.Println("✨ VALID EVENT: Received a FaceMatch event for a known human ID.")

		userName := gjsonResult.Get("params.events.0.data.alarmResult.faces.identify.candidate.reserve_field.name").String()
		deviceId := gjsonResult.Get("params.events.0.srcName").String()

		log.Printf("👤 User Name: %s", userName)
		log.Printf("📍 Device ID: %s", deviceId)

		go func() {
			log.Println("➡️ Sending data to internal authorization service...")

			dataToSend := map[string]string{
				"user_name": userName,
				"device_id": deviceId,
			}
			submissionObj, err := json.Marshal(dataToSend)
			if err != nil {
				log.Printf("❌ ERROR: Failed to marshal submission object! Reason: %v", err)
				return
			}
			submissionBytesArray := bytes.NewBuffer(submissionObj)

			tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			internalclient := &http.Client{Transport: tr}

			submissionUrl, err := url.JoinPath("http://localhost:"+os.Getenv("PORT"), "/api/v1/authorization/authorize-for-door-access")
			if err != nil {
				log.Printf("❌ ERROR: Failed to create submission URL! Reason: %v", err)
				return
			}

			httpReq, err := http.NewRequest("POST", submissionUrl, submissionBytesArray)
			if err != nil {
				log.Printf("❌ ERROR: Failed to create internal HTTP request! Reason: %v", err)
				return
			}

			httpReq.Header.Set("Content-Type", "application/json")
			httpReq.Header.Set("Internal-API-Key", internalkey.GetInternalAPIKey())

			resp, err := internalclient.Do(httpReq)
			if err != nil {
				log.Printf("❌ ERROR: Failed to send request to internal API! Reason: %v", err)
				return
			}
			defer resp.Body.Close()

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("❌ ERROR: Failed to read response body from internal API! Reason: %v", err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("⚠️ WARNING: Unexpected status code from internal API: %d", resp.StatusCode)
				log.Printf("Response Body: %s", string(bodyBytes))
			} else {
				log.Println("✅ SUCCESS: Internal API call completed successfully.")
				log.Printf("Final Response Status: %s", resp.Status)
				log.Printf("Response Body: %s", string(bodyBytes))
			}
		}()
		c.JSON(200, gin.H{"message": "successfully received event"})

	} else {
		c.JSON(200, gin.H{"message": "event received but not processed as a valid FaceMatch event"})
	}
}
