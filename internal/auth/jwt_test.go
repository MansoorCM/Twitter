package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	validSecret := "secret"
	validToken, _ := MakeJWT(userId, validSecret, time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserId  uuid.UUID
		wantErr     bool
	}{
		{name: "Valid Token",
			tokenString: validToken,
			tokenSecret: validSecret,
			wantUserId:  userId,
			wantErr:     false},
		{name: "Invalid Token",
			tokenString: "this is an invalid token",
			tokenSecret: validSecret,
			wantUserId:  uuid.Nil,
			wantErr:     true},
		{name: "Wrong Secret",
			tokenString: validToken,
			tokenSecret: "new secret",
			wantUserId:  uuid.Nil,
			wantErr:     true},
	}

	for _, tt := range tests {
		resUserId, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
		if (err != nil) != tt.wantErr {
			t.Errorf("validateJWT() - error %v, wanterr -%v", err, tt.wantErr)
		}
		if resUserId != tt.wantUserId {
			t.Errorf("validateJWT() - resUserId %v, wantUserId %v", resUserId, tt.wantUserId)
		}
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name       string
		header     http.Header
		wantHeader string
		wantErr    bool
	}{
		{name: "Valid token",
			header: http.Header{
				"Authorization": []string{"Bearer valid_token"}},
			wantHeader: "valid_token",
			wantErr:    false},
		{name: "Missing Authorization Header",
			header: http.Header{
				"Authorization": []string{}},
			wantHeader: "",
			wantErr:    true},
		{name: "Malformed Authorization header",
			header: http.Header{
				"Authorization": []string{"NotBearer invalid_token"}},
			wantHeader: "",
			wantErr:    true},
	}

	for _, tt := range tests {
		token, err := GetBearerToken(tt.header)
		if (err != nil) != tt.wantErr {
			t.Errorf("GetBearerToken(), err %v, wantErr %v", err, tt.wantErr)
		}
		if token != tt.wantHeader {
			t.Errorf("GetBearerToken(), gotToken %v, wantToken %v", token, tt.wantHeader)
		}
	}
}
