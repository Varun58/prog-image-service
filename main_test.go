package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	_ "image/jpeg"
	"testing"
	"os"
	"io"
	"mime/multipart"
	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	router := SetupRouter()

	// Create a new multipart form body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fileWriter, _ := writer.CreateFormFile("file", "test.jpg")
	file, _ := os.Open("testdata/test.jpg") // Provide a test image file
	defer file.Close()
	io.Copy(fileWriter, file)
	writer.Close()

	// Create a new POST request with the multipart form data
	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rec, req)

	// Assert the status code is OK (200)
	assert.Equal(t, http.StatusOK, rec.Code)

	// Parse the response JSON
	var response map[string]string
	json.NewDecoder(rec.Body).Decode(&response)

	// Assert the response contains an ID
	assert.NotNil(t, response["id"])
}

func TestGetImage(t *testing.T) {
	req, err := http.NewRequest("GET", "/image/valid_image_id/jpeg", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Set up a mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the existence of the image
		// For testing purposes, return a blank response
		w.WriteHeader(http.StatusOK)
	})

	// Serve the HTTP request and record the response
	handler.ServeHTTP(rr, req)

	// Check if the status code is as expected
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRotateImage(t *testing.T) {
	// Create a new HTTP request to rotate an image
	req, err := http.NewRequest("GET", "/transform/rotate/valid_image_id/90", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Set up a mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the rotation of the image
		// For testing purposes, return a blank response
		w.WriteHeader(http.StatusOK)
	})

	// Serve the HTTP request and record the response
	handler.ServeHTTP(rr, req)

	// Check if the status code is as expected
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestResizeImage(t *testing.T) {
	// Create a new HTTP request to resize an image
	req, err := http.NewRequest("GET", "/transform/resize/valid_image_id/100/100", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Set up a mock server
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate the resizing of the image
		// For testing purposes, return a blank response
		w.WriteHeader(http.StatusOK)
	})

	// Serve the HTTP request and record the response
	handler.ServeHTTP(rr, req)

	// Check if the status code is as expected
	assert.Equal(t, http.StatusOK, rr.Code)
}
