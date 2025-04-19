package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"uniproj/db"
)

func TestRegisterHandler(t *testing.T) {
	db.InitDB()
	defer db.CloseDB()

	reqBody, _ := json.Marshal(map[string]string{
		"email":    "test3@example.com",
		"password": "StrongPass123!",
	})
	req, _ := http.NewRequest("POST", "/auth/register", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	RegisterHandler(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Ожидался код статуса %d, получен %d", http.StatusCreated, rr.Code)
	}

	var response map[string]string
	json.NewDecoder(rr.Body).Decode(&response)
	if response["message"] != "User registered successfully" {
		t.Errorf("Ожидалось сообщение 'User registered successfully', получено '%s'", response["message"])
	}
}
