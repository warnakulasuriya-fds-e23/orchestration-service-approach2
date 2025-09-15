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
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils"
)

type EventReceiveController struct{}

func (erc *EventReceiveController) ReceiveFaceMatchEvent(c *gin.Context) {
	BodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error while reading request body": err.Error()})
		return
	}

	gjsonResult := gjson.ParseBytes(BodyBytes)
	// Check if the event is a FaceMatch event with a recognized human_id
	if gjsonResult.Get("method").String() == "OnEventNotify" &&
		gjsonResult.Get("params.ability").String() == "event_face_match" &&
		gjsonResult.Get("params.events.0.eventType").String() == "131659" &&
		gjsonResult.Get("params.events.0.data.alarmResult.faces.identify.candidate.human_id").String() != "-1" {
		log.Println("Received a valid FaceMatch event with a recognized human_id.")

		HcpAddress := os.Getenv("HCP_ADDRESS")
		HcpPersonInfoUrl, err := url.JoinPath("https://", HcpAddress, "/artemis/api/resource/v1/person/personId/personInfo")
		if err != nil {
			c.JSON(500, gin.H{"error while creating HCP person info URL": err.Error()})
			return
		}
		personalInfoReqObj := models.PersonalInfoRequestObj{
			PersonId:   gjsonResult.Get("params.events.0.data.alarmResult.faces.identify.candidate.human_id").String(),
			AppendInfo: []int{6},
		}
		jsonReqBody, err := json.Marshal(personalInfoReqObj)
		if err != nil {
			c.JSON(500, gin.H{"error while marshalling personal info request object to JSON": err.Error()})
			return
		}
		req, err := http.NewRequest("POST", HcpPersonInfoUrl, bytes.NewBuffer(jsonReqBody))
		if err != nil {
			c.JSON(500, gin.H{"error while creating new request to internal authorization endpoint": err.Error()})
			return
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-CA-Key", os.Getenv("HCP_OPENAPI_USER_KEY"))
		requestURLPath := "/artemis/api/resource/v1/person/personId/personInfo"
		generatedSignature, err := utils.GenerateSignatureForHcpOpenapi("POST", req.Header, requestURLPath, os.Getenv("HCP_OPENAPI_USER_SECRET"))
		if err != nil {
			c.JSON(500, gin.H{"error while generating signature for HCP OpenAPI": err.Error()})
			return
		}
		req.Header.Set("X-CA-Signature", generatedSignature)
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		internalclient := &http.Client{Transport: tr}
		resp, err := internalclient.Do(req)
		if err != nil {
			c.JSON(500, gin.H{"error while sending request to internal authorization endpoint": err.Error()})
			return
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error while reading response body from internal authorization endpoint": err.Error()})
			return
		}
		if resp.StatusCode != 200 {
			c.JSON(500, gin.H{"error while getting response from internal authorization endpoint": resp.Status})
			return
		}

		gjsonSecondResult := gjson.ParseBytes(bodyBytes)
		userId := gjsonSecondResult.Get("data.remark").String()
		deviceId := gjsonResult.Get("params.events.0.srcName").String()
		go func() {
			dataToSend := map[string]string{
				"user_id":   userId,
				"device_id": deviceId,
			}
			submissionObj, err := json.Marshal(dataToSend)
			if err != nil {
				log.Printf("Error marshalling submission object: %v", err)
				return
			}
			submissionBytesArray := bytes.NewBuffer(submissionObj)

			tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			internalclient := &http.Client{Transport: tr}
			submissionUrl, err := url.JoinPath("http://localhost:"+os.Getenv("PORT"), "/api/v1/authorization/authorize-for-door-access")
			if err != nil {
				log.Printf("Error creating submission URL: %v", err)
				return
			}
			httpReq, err := http.NewRequest("POST", submissionUrl, submissionBytesArray)
			if err != nil {
				log.Printf("Error creating HTTP request: %v", err)
				return
			}

			httpReq.Header.Set("Content-Type", "application/json")
			resp, err := internalclient.Do(httpReq)
			if err != nil {
				log.Printf("Error sending HTTP request: %v", err)
				return
			}
			defer resp.Body.Close()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("Error reading response body: %v", err)
				return
			}
			if resp.StatusCode != http.StatusOK {
				log.Printf("Unexpected status code from internal API: %d", resp.StatusCode)
				log.Printf("Response body: %v", string(bodyBytes))
			} else {
				log.Printf("Successfully processed FaceMatch event. Response: %s", resp.Status)
				log.Printf("response body: %v", string(bodyBytes))
			}
		}()
		c.JSON(200, gin.H{"message": "success"})
	}
}
