package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no auth header included in the request")
	}

	keySplit := strings.Split(authHeader, " ")
	if len(keySplit) != 2 || keySplit[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return keySplit[1], nil
}
