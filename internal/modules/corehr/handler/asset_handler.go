package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
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
// @Summary Create a new asset
// @Description Create a new asset with the provided details
// @Tags Assets
// @Accept json
// @Produce json
// @Param asset body model.Asset true "Asset object"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets [post]
func (h *AssetHandler) CreateAsset(c *fiber.Ctx) error {
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

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Asset created successfully"})
}

// GetAssetByID handles GET /assets/:id
// @Summary Get asset by ID
// @Description Get detailed information about an asset
// @Tags Assets
// @Produce json
// @Param id path int true "Asset ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.Asset
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/assets/{id} [get]
func (h *AssetHandler) GetAssetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid asset ID"})
	}

	asset, err := h.service.GetAssetByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Asset not found"})
	}

	return c.Status(fiber.StatusOK).JSON(asset)
}

// GetAssetBySerialNumber handles GET /assets/serial/:serialNumber
// @Summary Get asset by serial number
// @Description Get detailed information about an asset by its serial number
// @Tags Assets
// @Produce json
// @Param serialNumber path string true "Serial Number"
// @Security ApiKeyAuth
// @Success 200 {object} model.Asset
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/assets/serial/{serialNumber} [get]
func (h *AssetHandler) GetAssetBySerialNumber(c *fiber.Ctx) error {
	serialNumber := c.Params("serialNumber")
	asset, err := h.service.GetAssetBySerialNumber(serialNumber)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Asset not found"})
	}

	return c.Status(fiber.StatusOK).JSON(asset)
}

// GetAssetByAssetTag handles GET /assets/tag/:assetTag
// @Summary Get asset by asset tag
// @Description Get detailed information about an asset by its asset tag
// @Tags Assets
// @Produce json
// @Param assetTag path string true "Asset Tag"
// @Security ApiKeyAuth
// @Success 200 {object} model.Asset
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/assets/tag/{assetTag} [get]
func (h *AssetHandler) GetAssetByAssetTag(c *fiber.Ctx) error {
	assetTag := c.Params("assetTag")
	asset, err := h.service.GetAssetByAssetTag(assetTag)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Asset not found"})
	}

	return c.Status(fiber.StatusOK).JSON(asset)
}

// GetAllAssets handles GET /assets
// @Summary Get all assets
// @Description Get a list of all assets with pagination
// @Tags Assets
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security ApiKeyAuth
// @Success 200 {array} model.Asset
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets [get]
func (h *AssetHandler) GetAllAssets(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")

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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assets)
}

// UpdateAsset handles PUT /assets/:id
// @Summary Update an asset
// @Description Update an existing asset's details
// @Tags Assets
// @Accept json
// @Produce json
// @Param id path int true "Asset ID"
// @Param asset body model.Asset true "Asset update object"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets/{id} [put]
func (h *AssetHandler) UpdateAsset(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid asset ID"})
	}

	var body model.Asset
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	body.ID = uint(id)

	if err := h.service.UpdateAsset(&body); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset updated successfully"})
}

// DeleteAsset handles DELETE /assets/:id
// @Summary Delete an asset
// @Description Delete an asset by its ID
// @Tags Assets
// @Produce json
// @Param id path int true "Asset ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets/{id} [delete]
func (h *AssetHandler) DeleteAsset(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid asset ID"})
	}

	if err := h.service.DeleteAsset(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset deleted successfully"})
}

// GetAssetsByStatus handles GET /assets/status/:status
// @Summary Get assets by status
// @Description Get a list of assets with a specific status
// @Tags Assets
// @Produce json
// @Param status path string true "Status (available, assigned, maintenance, retired)"
// @Security ApiKeyAuth
// @Success 200 {array} model.Asset
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets/status/{status} [get]
func (h *AssetHandler) GetAssetsByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	assets, err := h.service.GetAssetsByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assets)
}

// GetAssetsByType handles GET /assets/type/:type
// @Summary Get assets by type
// @Description Get a list of assets of a specific type
// @Tags Assets
// @Produce json
// @Param type path string true "Asset Type (laptop, phone, tablet, etc.)"
// @Security ApiKeyAuth
// @Success 200 {array} model.Asset
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets/type/{type} [get]
func (h *AssetHandler) GetAssetsByType(c *fiber.Ctx) error {
	assetType := c.Params("type")
	assets, err := h.service.GetAssetsByType(assetType)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(assets)
}