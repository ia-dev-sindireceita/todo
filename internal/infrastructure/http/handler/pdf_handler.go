package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/usecases"
)

// PDFHandler handles HTTP requests for PDF export
type PDFHandler struct {
	exportTasksPDF usecases.ExportTasksPDFUseCaseInterface
}

// NewPDFHandler creates a new PDFHandler
func NewPDFHandler(exportTasksPDF usecases.ExportTasksPDFUseCaseInterface) *PDFHandler {
	return &PDFHandler{
		exportTasksPDF: exportTasksPDF,
	}
}

// ExportTasks handles GET /api/tasks/export/pdf
func (h *PDFHandler) ExportTasks(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID := r.Context().Value("userID").(string)

	// Generate PDF
	pdfBytes, err := h.exportTasksPDF.Execute(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	// Set headers for PDF download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=tarefas_%s.pdf", time.Now().Format("20060102_150405")))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfBytes)))

	// Write PDF to response
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBytes)
}
