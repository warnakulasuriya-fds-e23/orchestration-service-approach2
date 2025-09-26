package tokenstorage

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

type tokenResponseObject struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	ExpiresIn        int    `json:"expires_in"`
	Scope            string `json:"scope,omitempty"`
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// TokenStorage is to be used to manage the token that will be used to call the IDP
type TokenStorage struct {
	accessToken string
	expiryTime  time.Time
}

var (
	instance *TokenStorage
	once     sync.Once
)

func GetTokenStorage() (tokenStorage *TokenStorage) {
	once.Do(func() {
		instance = &TokenStorage{accessToken: "", expiryTime: time.Now()}
		_, err := instance.GetAccessToken()
		if err != nil {
			log.Fatalf("failed to initialize token storage: %v", err)
			tokenStorage = nil
			return
		}
	})
	tokenStorage = instance
	return
}

func (tokenStorage *TokenStorage) GetAccessToken() (token string, err error) {
	if tokenStorage.accessToken == "" || tokenStorage.expiryTime.Equal(time.Now()) || tokenStorage.expiryTime.Before(time.Now().Add(5*time.Second)) {

		tokenEndpoint, errTokenEndpoint := url.JoinPath(os.Getenv("IDP_BASE_URL"), "/oauth2/token")
		if errTokenEndpoint != nil {
			return "", fmt.Errorf("error while creating token endpoint URL: %w", errTokenEndpoint)
		}
		data := url.Values{}
		data.Set("grant_type", "client_credentials")
		data.Set("scope", "internal_user_mgt_view internal_user_mgt_list")
		requestBody := bytes.NewBufferString(data.Encode())
		req, errNewReq := http.NewRequest("POST", tokenEndpoint, requestBody)
		if errNewReq != nil {
			token = ""
			err = fmt.Errorf("error while creating a post request for the tokenEndpoint : %w", errNewReq)
			return
		}
		consumerKey := os.Getenv("CLIENT_ID")
		consumerSecret := os.Getenv("CLIENT_SECRET")
		authHeadervalue := base64.StdEncoding.EncodeToString([]byte(consumerKey + ":" + consumerSecret))
		req.Header.Add("Authorization", "Basic "+authHeadervalue)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		internalclient := &http.Client{Transport: tr}
		res, errReqSend := internalclient.Do(req)
		if errReqSend != nil {
			token = ""
			err = fmt.Errorf("error while sending or recieving post request : %w", errReqSend)
			return
		}
		defer res.Body.Close()
		bodybytes, errReadAll := io.ReadAll(res.Body)
		if errReadAll != nil {
			token = ""
			err = fmt.Errorf("error while reading bytes of response body : %w", errReadAll)
			return
		}
		var resObj tokenResponseObject
		errUnMarshal := json.Unmarshal(bodybytes, &resObj)
		if errUnMarshal != nil {
			token = ""
			err = fmt.Errorf("error while running json unmarshal for the read bytes of the response body : %w", errUnMarshal)
			return
		}

		tokenStorage.expiryTime = time.Now().Add(time.Duration(resObj.ExpiresIn) * time.Second)
		tokenStorage.accessToken = resObj.AccessToken
		log.Println("obtained new access token from : ", tokenEndpoint)

	}
	token = tokenStorage.accessToken
	err = nil
	return
}
