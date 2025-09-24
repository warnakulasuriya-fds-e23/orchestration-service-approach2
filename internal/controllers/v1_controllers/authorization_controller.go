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
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/utils/tokenstorage"
)

type AuthorizationController struct{}

func (ac *AuthorizationController) AuthorizeUserForDoorAccess(c *gin.Context) {
	var reqBody models.SubmissionForAuthorization
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error cannot bind json to incoming data from hikcentral format": err.Error()})
		return
	}

	userName, err := reqBody.GetUserName()
	if err != nil {
		c.JSON(500, gin.H{"error while getting user ID from request body": err.Error()})
		return
	}

	deviceId, err := reqBody.GetDeviceId()
	if err != nil {
		c.JSON(500, gin.H{"error while getting device ID from request body": err.Error()})
		return
	}

	idpAddress := os.Getenv("IDP_BASE_URL")
	scimCallUrl, err := url.JoinPath(idpAddress, "/scim2/Users")
	scimCallUrl = scimCallUrl + "?filter=userName+Co+\"" + userName + "\""
	if err != nil {
		c.JSON(500, gin.H{"error while joining URL path": err.Error()})
		return
	}
	newRequest, err := http.NewRequest("GET", scimCallUrl, nil)
	if err != nil {
		c.JSON(500, gin.H{"error while creating new request": err.Error()})
		return
	}
	newRequest.Header.Set("Accept", "application/scim+json")
	token, err := tokenstorage.GetTokenStorage().GetAccessToken()
	if err != nil {
		c.JSON(500, gin.H{"error while getting access token": err.Error()})
		return
	}
	newRequest.Header.Set("Authorization", "Bearer "+token)

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}
	resp, err := internalclient.Do(newRequest)
	if err != nil {
		c.JSON(500, gin.H{"error while making request to IDP": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to authorize user", "details": resp.Status, "idpresponse": resp.Body})
		resBodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error while reading response body": err.Error()})
			return
		}
		log.Println("Response from IDP:", string(resBodyBytes))
		return
	}
	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error while reading response body": err.Error()})
		return
	}
	log.Println(deviceId)
	resBody := gjson.ParseBytes(resBodyBytes)
	lengthOfGroups := resBody.Get("Resources.0.groups.#").Int()

	var groupNames []string
	for i := int64(0); i < lengthOfGroups; i++ {
		groupNames = append(groupNames, resBody.Get("Resources.0.groups."+strconv.Itoa(int(i))+".display").String())
	}
	accessGranted, err := utils.GroupBasedAuthorization(deviceId, groupNames)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to authorize user for device access"})
		return
	}
	if !accessGranted {
		c.JSON(403, gin.H{"error": "User is not authorized to access this device"})
		return
	}
	requirementManager := utils.GetRequirementsManager()
	doorId, err := requirementManager.GetDoorId(deviceId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	unlocked, err := utils.UnlockDoor(doorId)
	if err != nil || !unlocked {
		c.JSON(500, gin.H{"error": "Failed to unlock the door"})
		return
	}
	c.JSON(200, gin.H{"message": "User is authorized to access this device"})
}
