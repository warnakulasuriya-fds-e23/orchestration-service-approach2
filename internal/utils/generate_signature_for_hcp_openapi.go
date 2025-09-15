package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
)

func GenerateSignatureForHcpOpenapi(httpMethod string, header http.Header, requestPath, userSecret string) (string, error) {

	stringToSign := fmt.Sprintf("%s\n%s\n%s\n%s", httpMethod, header.Get("Accept"), header.Get("Content-Type"), requestPath)
	h := hmac.New(sha256.New, []byte(userSecret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}
