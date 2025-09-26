package v1_controllers

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/authorizationscache"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/internalkey"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/tokenstorage"
)

type AuthorizationController struct{}

func (ac *AuthorizationController) AuthorizeUserForDoorAccess(c *gin.Context) {
	log.Println("--- 🔑 AUTHORIZATION REQUEST RECEIVED 🔑 ---")

	if c.Request.Header.Get("Internal-API-Key") != internalkey.GetInternalAPIKey() {
		log.Println("❌ ERROR: Forbidden. Invalid Internal-API-Key provided.")
		c.JSON(403, gin.H{"error": "Forbidden"})
		return
	}

	var reqBody models.SubmissionForAuthorization
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("❌ ERROR: Failed to bind JSON. Reason: %v", err)
		c.JSON(400, gin.H{"error cannot bind json to incoming data from hikcentral format": err.Error()})
		return
	}

	userName, err := reqBody.GetUserName()
	if err != nil {
		log.Printf("❌ ERROR: Failed to get user name from request body. Reason: %v", err)
		c.JSON(500, gin.H{"error while getting user ID from request body": err.Error()})
		return
	}
	log.Printf("👤 User Name received: %s", userName)

	deviceId, err := reqBody.GetDeviceId()
	if err != nil {
		log.Printf("❌ ERROR: Failed to get device ID from request body. Reason: %v", err)
		c.JSON(500, gin.H{"error while getting device ID from request body": err.Error()})
		return
	}
	log.Printf("🚪 Device ID received: %s", deviceId)

	requirementManager := utils.GetRequirementsManager()
	doorId, err := requirementManager.GetDoorId(deviceId)
	if err != nil {
		log.Printf("❌ ERROR: Failed to get door ID for device '%s'. Reason: %v", deviceId, err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// check AuthorizationsCache
	cachedAuth := authorizationscache.GetAuthorizationsCacheInstance()
	if cachedAuth.IsAuthorized(userName, doorId) {
		log.Printf("✅ CACHE HIT: User '%s' is already authorized for door '%s'.", userName, doorId)

		unlocked, err := utils.UnlockDoor(doorId)
		if err != nil || !unlocked {
			log.Printf("❌ ERROR: Failed to unlock door '%s'. Reason: %v", doorId, err)
			c.JSON(500, gin.H{"error": "Failed to unlock the door"})
			return
		}
		log.Println("🔓 SUCCESS: Door unlocked!")

		log.Println("--- 🔓 AUTHORIZATION PROCESS COMPLETE 🔓 ---")
		c.JSON(200, gin.H{"message": "User is authorized to access this device"})
	}

	log.Printf("🚪 Door ID '%s' found for device '%s'.", doorId, deviceId)
	idpAddress := os.Getenv("IDP_BASE_URL")
	scimCallUrl, err := url.JoinPath(idpAddress, "/scim2/Users")
	scimCallUrl = scimCallUrl + "?filter=userName+Co+\"" + userName + "\""
	if err != nil {
		log.Printf("❌ ERROR: Failed to join URL path. Reason: %v", err)
		c.JSON(500, gin.H{"error while joining URL path": err.Error()})
		return
	}
	log.Printf("➡️ Sending request to IDP at: %s", scimCallUrl)

	newRequest, err := http.NewRequest("GET", scimCallUrl, nil)
	if err != nil {
		log.Printf("❌ ERROR: Failed to create new HTTP request. Reason: %v", err)
		c.JSON(500, gin.H{"error while creating new request": err.Error()})
		return
	}
	newRequest.Header.Set("Accept", "application/scim+json")

	token, err := tokenstorage.GetTokenStorage().GetAccessToken()
	if err != nil {
		log.Printf("❌ ERROR: Failed to get access token. Reason: %v", err)
		c.JSON(500, gin.H{"error while getting access token": err.Error()})
		return
	}
	newRequest.Header.Set("Authorization", "Bearer "+token)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}
	resp, err := internalclient.Do(newRequest)
	if err != nil {
		log.Printf("❌ ERROR: Failed to make request to IDP. Reason: %v", err)
		c.JSON(500, gin.H{"error while making request to IDP": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("⚠️ WARNING: IDP returned an unexpected status code: %d", resp.StatusCode)
		resBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("❌ ERROR: Failed to read IDP response body. Reason: %v", err)
			c.JSON(500, gin.H{"error while reading response body": err.Error()})
			return
		}
		log.Printf("IDP Response Body: %s", string(resBodyBytes))
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to authorize user", "details": resp.Status, "idpresponse": string(resBodyBytes)})
		return
	}

	log.Println("✅ IDP responded with 200 OK. Reading response body...")
	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("❌ ERROR: Failed to read response body from IDP. Reason: %v", err)
		c.JSON(500, gin.H{"error while reading response body": err.Error()})
		return
	}
	resBody := gjson.ParseBytes(resBodyBytes)
	log.Printf("IDP response body: %s", string(resBodyBytes))

	lengthOfGroups := resBody.Get("Resources.0.groups.#").Int()
	var groupNames []string
	for i := int64(0); i < lengthOfGroups; i++ {
		groupNames = append(groupNames, resBody.Get("Resources.0.groups."+strconv.Itoa(int(i))+".display").String())
	}
	log.Printf("👥 User belongs to the following groups: %v", groupNames)

	accessGranted, err := utils.GroupBasedAuthorization(deviceId, groupNames)
	if err != nil {
		log.Printf("❌ ERROR: Group-based authorization failed. Reason: %v", err)
		c.JSON(500, gin.H{"error": "Failed to authorize user for device access"})
		return
	}

	if !accessGranted {
		log.Println("⛔ ACCESS DENIED: User is not authorized for this device based on group membership.")
		c.JSON(403, gin.H{"error": "User is not authorized to access this device"})
		return
	}

	log.Println("✅ ACCESS GRANTED: User is authorized for this device.")

	// cache the authorization
	cachedAuth.SetAuthorization(userName, doorId)

	unlocked, err := utils.UnlockDoor(doorId)
	if err != nil || !unlocked {
		log.Printf("❌ ERROR: Failed to unlock door '%s'. Reason: %v", doorId, err)
		c.JSON(500, gin.H{"error": "Failed to unlock the door"})
		return
	}
	log.Println("🔓 SUCCESS: Door unlocked!")

	log.Println("--- 🔓 AUTHORIZATION PROCESS COMPLETE 🔓 ---")
	c.JSON(200, gin.H{"message": "User is authorized to access this device"})
}
