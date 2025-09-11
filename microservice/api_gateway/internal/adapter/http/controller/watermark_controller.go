package controller

import (
	"encoding/json"
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

	var req request.ApplyWatermarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad requests", http.StatusBadRequest)
		return
	}

	res, err := h.service.ApplyWatermark(r.Context(), req)
	if err != nil {
		http.Error(w, "gateway error", http.StatusBadGateway)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
