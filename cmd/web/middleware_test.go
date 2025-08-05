package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"snippetbox.alexedwards.net/internal/assert"
	"testing"
)

func TestSecureHeaders(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	// create dummy handler that simulate the handler that comes after secureHeaders middleware
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	// wrap the handler with secureHeaders middleware
	secureHeaders(next).ServeHTTP(rr, req)

	res := rr.Result()

	expectedValue := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, res.Header.Get("Content-Security-Policy"), expectedValue)

	expectedValue = "origin-when-cross-origin"
	assert.Equal(t, res.Header.Get("referrer-policy"), expectedValue)

	expectedValue = "nosniff"
	assert.Equal(t, res.Header.Get("X-Content-Type-Options"), expectedValue)

	expectedValue = "deny"
	assert.Equal(t, res.Header.Get("X-Frame-Options"), expectedValue)

	expectedValue = "0"
	assert.Equal(t, res.Header.Get("X-XSS-Protection"), expectedValue)

	assert.Equal(t, res.StatusCode, http.StatusOK)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	bytes.TrimSpace(body)
	assert.Equal(t, string(body), "ok")
}
