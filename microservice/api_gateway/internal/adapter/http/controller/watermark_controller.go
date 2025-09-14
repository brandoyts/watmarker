package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/model/request"
	"github.com/brandoyts/watmarker/microservice/api_gateway/internal/core/port"
)

type WatermarkController struct {
	service port.WatermarkService
}

func NewWatermarkController(s port.WatermarkService) *WatermarkController {
	return &WatermarkController{service: s}
}

func (h *WatermarkController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	watermarkText := r.FormValue("watermark_text")
	if watermarkText == "" {
		http.Error(w, "missing watermark_text field", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file_data")
	if err != nil {
		http.Error(w, "Error retrieving file from form", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		log.Printf("Error reading file content: %v", err)
		return
	}

	res, err := h.service.ApplyWatermark(r.Context(), request.ApplyWatermarkRequest{
		WatermarkText: watermarkText,
		FileData:      fileBytes,
	})
	if err != nil {
		http.Error(w, "failed to process the image", http.StatusInternalServerError)
		log.Printf("failed to process the image: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
