package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// Dummy reference to satisfy swag and lint
var _ = model.OrganizationNode{}

// OrganizationHandler handles HTTP requests for organizational structure operations
type OrganizationHandler struct {
	employeeService service.EmployeeService
}

// NewOrganizationHandler creates a new instance of OrganizationHandler
func NewOrganizationHandler(employeeService service.EmployeeService) *OrganizationHandler {
	return &OrganizationHandler{
		employeeService: employeeService,
	}
}

// GetTree handles GET /organization/tree
// @Summary Get organization hierarchy tree
// @Description Retrieve the full company structure in a hierarchical JSON format
// @Tags Organization
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} model.OrganizationNode
// @Failure 500 {object} map[string]string
// @Router /api/v1/organization/tree [get]
func (h *OrganizationHandler) GetTree(c *fiber.Ctx) error {
	tree, err := h.employeeService.GetOrganizationTree()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tree)
}
