package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockExportPDFUseCase struct {
	pdfBytes []byte
	err      error
}

func (m *MockExportPDFUseCase) Execute(ctx context.Context, ownerID string) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.pdfBytes, nil
}

func TestPDFHandler_ExportTasks(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockPDFBytes   []byte
		mockError      error
		expectedStatus int
		checkHeaders   bool
	}{
		{
			name:           "Successfully export tasks to PDF",
			userID:         "user-1",
			mockPDFBytes:   []byte("%PDF-1.4 test content"),
			mockError:      nil,
			expectedStatus: http.StatusOK,
			checkHeaders:   true,
		},
		{
			name:           "Error generating PDF",
			userID:         "user-1",
			mockPDFBytes:   nil,
			mockError:      errors.New("failed to generate PDF"),
			expectedStatus: http.StatusInternalServerError,
			checkHeaders:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUseCase := &MockExportPDFUseCase{
				pdfBytes: tt.mockPDFBytes,
				err:      tt.mockError,
			}

			handler := NewPDFHandler(mockUseCase)

			req := httptest.NewRequest(http.MethodGet, "/api/tasks/export/pdf", nil)
			ctx := context.WithValue(req.Context(), "userID", tt.userID)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			handler.ExportTasks(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.checkHeaders {
				contentType := w.Header().Get("Content-Type")
				if contentType != "application/pdf" {
					t.Errorf("Expected Content-Type application/pdf, got %s", contentType)
				}

				contentDisposition := w.Header().Get("Content-Disposition")
				if contentDisposition == "" {
					t.Error("Expected Content-Disposition header to be set")
				}

				if !bytes.Equal(w.Body.Bytes(), tt.mockPDFBytes) {
					t.Error("Response body does not match expected PDF bytes")
				}
			}
		})
	}
}
