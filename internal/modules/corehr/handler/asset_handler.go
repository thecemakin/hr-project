package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
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

	if err := h.service.CreateAsset(1, asset); err != nil {
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

func (h *AssetHandler) GetAssetBySerialNumber(c *fiber.Ctx) error {
	serialNumber := c.Params("serialNumber")
	filter := repository.AssetFilter{SerialNumber: serialNumber}
	assets, err := h.service.GetAllAssets(1, 0, filter)
	if err != nil || len(assets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Asset not found"})
	}
	return c.Status(fiber.StatusOK).JSON(assets[0])
}

func (h *AssetHandler) GetAssetByAssetTag(c *fiber.Ctx) error {
	assetTag := c.Params("assetTag")
	filter := repository.AssetFilter{AssetTag: assetTag}
	assets, err := h.service.GetAllAssets(1, 0, filter)
	if err != nil || len(assets) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Asset not found"})
	}
	return c.Status(fiber.StatusOK).JSON(assets[0])
}

// GetAllAssets handles GET /assets
// @Summary Get assets with filtering
// @Description Get a list of assets with pagination and optional filters (serial_number, asset_tag, status, type, search)
// @Tags Assets
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param serial_number query string false "Filter by serial number"
// @Param asset_tag query string false "Filter by asset tag"
// @Param status query string false "Filter by status"
// @Param type query string false "Filter by type"
// @Param search query string false "Search by brand or model"
// @Security ApiKeyAuth
// @Success 200 {array} model.Asset
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/assets [get]
func (h *AssetHandler) GetAllAssets(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	serialNumber := c.Query("serial_number")
	assetTag := c.Query("asset_tag")
	status := c.Query("status")
	assetType := c.Query("type")
	search := c.Query("search")

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	filter := repository.AssetFilter{
		SerialNumber: serialNumber,
		AssetTag:     assetTag,
		Status:       status,
		Type:         assetType,
		Search:       search,
	}

	assets, err := h.service.GetAllAssets(limit, offset, filter)
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

	if err := h.service.UpdateAsset(1, &body); err != nil {
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

	if err := h.service.DeleteAsset(1, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Asset deleted successfully"})
}
