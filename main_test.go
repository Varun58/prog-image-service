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