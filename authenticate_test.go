package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dhable/mini-auth/app"
	"github.com/dhable/mini-auth/models"
)

func TestAuthenticate(t *testing.T) {
	app, err := app.NewApp()
	if err != nil {
		t.Fatalf("failed to create application instance: %s", err)
	}

	t.Run("bad request", func(t *testing.T) {
		body := `{"something": "invalid"}`
		req, err := http.NewRequest("POST", "/authenicate", strings.NewReader(body))
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.AuthenticateUser)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusBadRequest {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("user doesn't exist", func(t *testing.T) {
		body := `{
			"email": "nobody@danhable.com",
			"password": "password1!"
		}`
		req, err := http.NewRequest("POST", "/authenicate", strings.NewReader(body))
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.AuthenticateUser)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusUnauthorized {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("user has incorrect password", func(t *testing.T) {
		body := `{
			"email": "dan@danhable.com",
			"password": "ssssh"
		}`
		req, err := http.NewRequest("POST", "/authenicate", strings.NewReader(body))
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.AuthenticateUser)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusUnauthorized {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusUnauthorized)
		}
	})

	t.Run("correct credentials", func(t *testing.T) {
		body := `{
			"email": "dan@danhable.com",
			"password": "password1!"
		}`
		req, err := http.NewRequest("POST", "/authenicate", strings.NewReader(body))
		if err != nil {
			t.Fatalf("failed to create http request: %s", err)
		}

		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(app.AuthenticateUser)

		handler.ServeHTTP(recorder, req)

		status := recorder.Code
		if status != http.StatusOK {
			t.Errorf("incorrect status code: got %v want %v", status, http.StatusOK)
		}

		var resp models.AuthResponse
		err = json.Unmarshal(recorder.Body.Bytes(), &resp)
		if err != nil {
			t.Errorf("failed to parse body: %s", err)
		}

		if resp.Jwt == "" {
			t.Error("expected Jwt value")
		}
	})
}
