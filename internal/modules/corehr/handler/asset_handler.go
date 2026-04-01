package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// AssetHandler handles HTTP requests for asset operations
type AssetHandler struct {
	service service.AssetService
}

// NewAssetHandler creates a new instance of AssetHandler
func NewAssetHandler(service service.AssetService) *AssetHandler {
	return &AssetHandler{
		service: service,
	}
}

// CreateAsset handles POST /assets
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
	var body struct {
		SerialNumber    string  `json:"serial_number"`
		Name            string  `json:"name"`
		Description     string  `json:"description"`
		AssetTag        string  `json:"asset_tag"`
		Type            string  `json:"type"`
		Brand           string  `json:"brand"`
		Model           string  `json:"model"`
		PurchaseDate    *string `json:"purchase_date"`
		PurchasePrice   float64 `json:"purchase_price"`
		Status          string  `json:"status"`
		Condition       string  `json:"condition"`
		WarrantyExpires *string `json:"warranty_expires"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	asset := &model.Asset{
		SerialNumber:    body.SerialNumber,
		Name:            body.Name,
		Description:     body.Description,
		AssetTag:        body.AssetTag,
		Type:            body.Type,
		Brand:           body.Brand,
		Model:           body.Model,
		PurchasePrice:   body.PurchasePrice,
		Status:          body.Status,
		Condition:       body.Condition,
	}

	if err := h.service.CreateAsset(asset); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Asset created successfully"})
}

// GetAssetByID handles GET /assets/:id
func (h *AssetHandler) GetAssetByID(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	asset, err := h.service.GetAssetByID(uint(id))
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Asset not found")
		return
	}

	respondWithJSON(w, http.StatusOK, asset)
}

// GetAssetBySerialNumber handles GET /assets/serial/:serialNumber
func (h *AssetHandler) GetAssetBySerialNumber(w http.ResponseWriter, r *http.Request) {
	serialNumber := chi.URLParam(r, "serialNumber")
	asset, err := h.service.GetAssetBySerialNumber(serialNumber)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Asset not found")
		return
	}

	respondWithJSON(w, http.StatusOK, asset)
}

// GetAssetByAssetTag handles GET /assets/tag/:assetTag
func (h *AssetHandler) GetAssetByAssetTag(w http.ResponseWriter, r *http.Request) {
	assetTag := chi.URLParam(r, "assetTag")
	asset, err := h.service.GetAssetByAssetTag(assetTag)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Asset not found")
		return
	}

	respondWithJSON(w, http.StatusOK, asset)
}

// GetAllAssets handles GET /assets
func (h *AssetHandler) GetAllAssets(w http.ResponseWriter, r *http.Request) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	assets, err := h.service.GetAllAssets(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assets)
}

// UpdateAsset handles PATCH /assets/:id
func (h *AssetHandler) UpdateAsset(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	var body model.Asset
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	body.ID = uint(id)

	if err := h.service.UpdateAsset(&body); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset updated successfully"})
}

// DeleteAsset handles DELETE /assets/:id
func (h *AssetHandler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid asset ID")
		return
	}

	if err := h.service.DeleteAsset(uint(id)); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Asset deleted successfully"})
}

// GetAssetsByStatus handles GET /assets/status/:status
func (h *AssetHandler) GetAssetsByStatus(w http.ResponseWriter, r *http.Request) {
	status := chi.URLParam(r, "status")
	assets, err := h.service.GetAssetsByStatus(status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assets)
}

// GetAssetsByType handles GET /assets/type/:type
func (h *AssetHandler) GetAssetsByType(w http.ResponseWriter, r *http.Request) {
	assetType := chi.URLParam(r, "type")
	assets, err := h.service.GetAssetsByType(assetType)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, assets)
}