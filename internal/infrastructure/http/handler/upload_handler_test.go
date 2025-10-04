package handler

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUploadImage_Success(t *testing.T) {
	// Create a temporary upload directory for testing
	tempDir := t.TempDir()

	handler := NewUploadHandler(tempDir)

	// Create a test image file with valid JPEG header
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.jpg")
	// JPEG magic bytes (SOI marker)
	jpegHeader := []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46}
	part.Write(jpegHeader)
	part.Write(make([]byte, 1000)) // Add some padding
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	handler.UploadImage(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	// Check that response contains the image path
	if !strings.Contains(w.Body.String(), "/uploads/images/") {
		t.Errorf("Response should contain image path: %s", w.Body.String())
	}
}

func TestUploadImage_FileTooLarge(t *testing.T) {
	tempDir := t.TempDir()
	handler := NewUploadHandler(tempDir)

	// Create a file larger than 10MB
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "large.jpg")
	largeContent := make([]byte, 11*1024*1024) // 11MB
	part.Write(largeContent)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	handler.UploadImage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestUploadImage_InvalidFileType(t *testing.T) {
	tempDir := t.TempDir()
	handler := NewUploadHandler(tempDir)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("image", "test.txt")
	part.Write([]byte("not an image"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload/image", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	handler.UploadImage(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
