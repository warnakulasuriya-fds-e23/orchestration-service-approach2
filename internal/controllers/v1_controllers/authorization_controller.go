package v1_controllers

import (
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/warnakulasuriya-fds-e23/orchestration-service-approach2/internal/models"
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

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		c.JSON(500, gin.H{"error while decoding response from IDP": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, res)
}
