package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dhable/mini-auth/app"
	"github.com/dhable/mini-auth/models"
)

func TestProfile(t *testing.T) {
	app, err := app.NewApp()
	if err != nil {
		t.Fatalf("failed to create application instance: %s", err)
	}

	t.Run("missing header", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/profile", nil)
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.CurrentUserProfile)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusUnauthorized {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("invalid authorization value", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/profile", nil)
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		req.Header.Add("Authorization", "garbage values")

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.CurrentUserProfile)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusUnauthorized {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("unknown jwt", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/profile", nil)
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		jwt := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c`
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.CurrentUserProfile)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusUnauthorized {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("valid user jwt", func(t *testing.T) {
		// Obtains a JWT using the authenticate endpoint
		body := `{
			"email": "dan@danhable.com",
			"password": "password1!"
		}`
		req, err := http.NewRequest("POST", "/authenicate", strings.NewReader(body))
		if err != nil {
			t.Fatalf("failed to create authenicate http request: %s", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.AuthenticateUser)

		handler.ServeHTTP(recorder, req)

		var authResp models.AuthResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &authResp)
		if err != nil {
			t.Errorf("failed to parse body: %s", err)
		}

		// Gets profile with that JWT
		req, err = http.NewRequest("GET", "/profile", nil)
		if err != nil {
			t.Fatalf("failed to create profile http request: %s", err)
		}

		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authResp.Jwt))

		recorder = httptest.NewRecorder()
		handler = http.HandlerFunc(app.CurrentUserProfile)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusOK {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusOK)
		}

		var profileResp models.UserProfile
		err = json.Unmarshal(recorder.Body.Bytes(), &profileResp)
		if err != nil {
			t.Errorf("failed to parse profile body: %s", err)
		}

		expectedEmail := "dan@danhable.com"
		if profileResp.Email != expectedEmail {
			t.Errorf("incorrect profile email: got %s want %s", profileResp.Email, expectedEmail)
		}

		expectedName := "Dan Hable"
		if profileResp.Name != expectedName {
			t.Errorf("incorrect profile name: got %s want %s", profileResp.Name, expectedName)
		}

		expectedLocation := "Chicago, IL"
		if profileResp.Location != expectedLocation {
			t.Errorf("incorrect profile location: got %s want %s", profileResp.Location, expectedLocation)
		}
	})
}
