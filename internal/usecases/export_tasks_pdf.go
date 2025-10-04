package usecases

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ia-edev-sindireceita/todo/internal/domain/application"
	"github.com/ia-edev-sindireceita/todo/internal/domain/repository"
	"github.com/jung-kurt/gofpdf"
)

// ExportTasksPDFUseCase handles exporting tasks to PDF
type ExportTasksPDFUseCase struct {
	taskRepo repository.TaskRepository
}

// NewExportTasksPDFUseCase creates a new ExportTasksPDFUseCase
func NewExportTasksPDFUseCase(taskRepo repository.TaskRepository) *ExportTasksPDFUseCase {
	return &ExportTasksPDFUseCase{
		taskRepo: taskRepo,
	}
}

// Execute generates a PDF with all tasks for a user
func (uc *ExportTasksPDFUseCase) Execute(ctx context.Context, ownerID string) ([]byte, error) {
	// Get all tasks for the user
	tasks, err := uc.taskRepo.FindByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve tasks: %w", err)
	}

	// Create PDF with UTF-8 support
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Add UTF-8 font support
	tr := pdf.UnicodeTranslatorFromDescriptor("")
	pdf.AddPage()

	// Set title
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(190, 10, tr("Minhas Tarefas"), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Add generation date
	pdf.SetFont("Arial", "I", 10)
	pdf.CellFormat(190, 6, tr(fmt.Sprintf("Gerado em: %s", time.Now().Format("02/01/2006 15:04:05"))), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add tasks
	if len(tasks) == 0 {
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(190, 10, tr("Nenhuma tarefa encontrada."), "", 1, "L", false, 0, "")
	} else {
		for i, task := range tasks {
			// Task number and title
			pdf.SetFont("Arial", "B", 14)
			pdf.CellFormat(190, 8, tr(fmt.Sprintf("%d. %s", i+1, task.Title)), "", 1, "L", false, 0, "")
			pdf.Ln(2)

			// Status
			pdf.SetFont("Arial", "", 11)
			statusText := getStatusText(task.Status)
			pdf.CellFormat(190, 6, tr(fmt.Sprintf("Status: %s", statusText)), "", 1, "L", false, 0, "")

			// Description
			if task.Description != "" {
				pdf.SetFont("Arial", "", 11)
				pdf.MultiCell(190, 5, tr(fmt.Sprintf("Descricao: %s", task.Description)), "", "L", false)
			}

			// Image (if present)
			if task.ImagePath != "" {
				// Convert relative path to absolute path
				imagePath := strings.TrimPrefix(task.ImagePath, "/")

				// Check if file exists
				if _, err := os.Stat(imagePath); err == nil {
					// Get current Y position
					currentY := pdf.GetY()

					// Register image and get dimensions
					opt := gofpdf.ImageOptions{
						ImageType: getImageType(imagePath),
						ReadDpi:   true,
					}

					// Calculate image size (max 200x200px = ~70x70mm at 72dpi)
					imgWidth := 70.0  // 200px at 72dpi ≈ 70mm
					imgHeight := 70.0

					// Add image with size constraints
					pdf.ImageOptions(imagePath, 10, currentY+2, imgWidth, imgHeight, false, opt, 0, "")

					// Move Y position after image
					pdf.SetY(currentY + imgHeight + 4)
				}
			}

			// Created date
			pdf.SetFont("Arial", "I", 9)
			pdf.CellFormat(190, 5, tr(fmt.Sprintf("Criada em: %s", task.CreatedAt.Format("02/01/2006 15:04"))), "", 1, "L", false, 0, "")

			// Add spacing between tasks
			pdf.Ln(8)
		}
	}

	// Output PDF to buffer
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}

// getStatusText converts task status to Portuguese text
func getStatusText(status application.TaskStatus) string {
	switch status {
	case application.StatusPending:
		return "Pendente"
	case application.StatusInProgress:
		return "Em Progresso"
	case application.StatusCompleted:
		return "Concluida"
	default:
		return "Desconhecido"
	}
}

// getImageType returns the image type for gofpdf based on file extension
func getImageType(imagePath string) string {
	ext := strings.ToLower(filepath.Ext(imagePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "JPEG"
	case ".png":
		return "PNG"
	case ".gif":
		return "GIF"
	default:
		return "JPEG" // default fallback
	}
}
