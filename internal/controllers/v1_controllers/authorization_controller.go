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
)

type AuthorizationController struct{}

func (ac *AuthorizationController) AuthorizeUserForDoorAccess(c *gin.Context) {
	var reqBody models.IncomingDataFromHikCentral
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(400, gin.H{"error cannot bind json to incoming data from hikcentral format": err.Error()})
		return
	}

	userID, err := reqBody.GetUserId()
	if err != nil {
		c.JSON(500, gin.H{"error while getting user ID from request body": err.Error()})
		return
	}

	deviceId, err := reqBody.GetDeviceId()
	if err != nil {
		c.JSON(500, gin.H{"error while getting device ID from request body": err.Error()})
		return
	}

	idpAddress := os.Getenv("IDP_ADDRESS")
	scimCallUrl, err := url.JoinPath(idpAddress, "/scim2/Users/", userID)
	if err != nil {
		c.JSON(500, gin.H{"error while joining URL path": err.Error()})
		return
	}
	newRequest, err := http.NewRequest("GET", scimCallUrl, nil)
	if err != nil {
		c.JSON(500, gin.H{"error while creating new request": err.Error()})
		return
	}
	newRequest.SetBasicAuth(os.Getenv("IDP_USERNAME"), os.Getenv("IDP_PASSWORD"))
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	internalclient := &http.Client{Transport: tr}
	resp, err := internalclient.Do(newRequest)
	if err != nil {
		c.JSON(500, gin.H{"error while making request to IDP": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": "Failed to authorize user"})
		return
	}
	resBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error while reading response body": err.Error()})
		return
	}
	log.Println(deviceId)
	resBody := gjson.ParseBytes(resBodyBytes)
	lengthOfRoles := resBody.Get("roles.#").Int()
	var roles []models.WSO2IDPRoleObject
	for i := int64(0); i < lengthOfRoles; i++ {
		roles = append(roles, models.WSO2IDPRoleObject{
			Ref:             resBody.Get("roles." + strconv.Itoa(int(i)) + ".$ref").String(),
			AudienceDisplay: resBody.Get("roles." + strconv.Itoa(int(i)) + ".audienceDisplay").String(),
			AudienceType:    resBody.Get("roles." + strconv.Itoa(int(i)) + ".audienceType").String(),
			AudienceValue:   resBody.Get("roles." + strconv.Itoa(int(i)) + ".audienceValue").String(),
			Display:         resBody.Get("roles." + strconv.Itoa(int(i)) + ".display").String(),
			Value:           resBody.Get("roles." + strconv.Itoa(int(i)) + ".value").String(),
		})
	}

	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.GetRoleName())
	}
	accessGranted, err := utils.RoleBasedAuthorization(deviceId, roleNames)
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
