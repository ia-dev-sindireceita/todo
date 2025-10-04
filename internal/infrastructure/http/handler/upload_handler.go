package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	MaxFileSize = 10 * 1024 * 1024 // 10MB
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
}

// UploadHandler handles file uploads
type UploadHandler struct {
	uploadDir string
}

// NewUploadHandler creates a new UploadHandler
func NewUploadHandler(uploadDir string) *UploadHandler {
	return &UploadHandler{
		uploadDir: uploadDir,
	}
}

// SaveImage saves an uploaded image and returns the relative path
func (h *UploadHandler) SaveImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validate file size
	if header.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds 10MB limit")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExtensions[ext] {
		return "", fmt.Errorf("invalid file type. Only JPG, PNG, GIF, and WebP are allowed")
	}

	// Read first 512 bytes to detect MIME type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("error reading file")
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Validate MIME type
	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return "", fmt.Errorf("invalid file type: %s. Only images are allowed", mimeType)
	}

	// Generate unique filename using hash
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("error processing file")
	}
	file.Seek(0, 0) // Reset for copying to disk

	hash := hex.EncodeToString(hasher.Sum(nil))
	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s%s", timestamp, hash[:16], ext)

	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(h.uploadDir, 0755); err != nil {
		return "", fmt.Errorf("error creating upload directory")
	}

	// Save file to disk
	filepath := filepath.Join(h.uploadDir, filename)
	dst, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("error saving file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("error saving file")
	}

	// Return the relative path
	relativePath := "/uploads/images/" + filename
	return relativePath, nil
}

// UploadImage handles image upload with security validations (HTTP endpoint)
func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Limit request body size to prevent DoS
	r.Body = http.MaxBytesReader(w, r.Body, MaxFileSize)

	// Parse multipart form with max memory limit
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB in memory
		http.Error(w, "File too large or invalid form data", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path, err := h.SaveImage(file, header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"path": path,
	})
}
