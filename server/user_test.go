package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSignup(t *testing.T) {
	GetDB()
	router := setupRouter()

	w := httptest.NewRecorder()
	input := `{"Email":"dane@turbobuilt.com"}`
	req, _ := http.NewRequest("POST", "/api/v1.0/user", bytes.NewBufferString(input))
	req.Header.Add("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	fmt.Println("Response code", w.Code, w.Body.String())
	assert.Equal(t, 200, w.Code)
	// router.
	// assert.Equal(t, "", w.Body.String())
}
