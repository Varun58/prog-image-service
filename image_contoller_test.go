package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUploadImage(t *testing.T) {
	// Create a new HTTP request with a file upload
	req := httptest.NewRequest("POST", "/upload", nil)
	file := bytes.NewBufferString("test file content") // Replace with your test file content
	req.Header.Set("Content-Type", "multipart/form-data")
	req.Header.Set("Content-Disposition", `form-data; name="file"; filename="test.jpg"`)

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Call the handler function
	UploadImage(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetImage(t *testing.T) {
	// Create a new HTTP request to get an image
	req := httptest.NewRequest("GET", "/image/test-image-id.jpg", nil)

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Call the handler function
	GetImage(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Optionally, you can also check the response body or headers
}

func TestRotateImage(t *testing.T) {
	// Create a new HTTP request to rotate an image
	req := httptest.NewRequest("GET", "/transform/rotate/test-image-id/90", nil)

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Call the handler function
	RotateImage(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Optionally, you can also check the response body or headers
}

func TestResizeImage(t *testing.T) {
	// Create a new HTTP request to resize an image
	req := httptest.NewRequest("GET", "/transform/resize/test-image-id/100/100", nil)

	// Create a response recorder to record the response
	rec := httptest.NewRecorder()

	// Call the handler function
	ResizeImage(rec, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rec.Code)

	// Optionally, you can also check the response body or headers
}
