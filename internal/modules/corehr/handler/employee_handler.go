package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
	"github.com/thecemakin/hr-project/internal/modules/corehr/repository"
	"github.com/thecemakin/hr-project/internal/modules/corehr/service"
)

// EmployeeHandler handles HTTP requests for employee operations
type EmployeeHandler struct {
	service service.EmployeeService
}

// NewEmployeeHandler creates a new instance of EmployeeHandler
func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		service: service,
	}
}

// CreateEmployee handles POST /employees
// @Summary Create a new employee
// @Description Create a new employee with the provided details
// @Tags Employees
// @Accept json
// @Produce json
// @Param employee body model.Employee true "Employee object"
// @Security ApiKeyAuth
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees [post]
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var body struct {
		FirstName             string `json:"first_name"`
		LastName              string `json:"last_name"`
		Email                 string `json:"email"`
		Phone                 string `json:"phone"`
		DateOfBirth           *string `json:"date_of_birth"`
		Gender                string `json:"gender"`
		AddressLine1          string `json:"address_line1"`
		AddressLine2          string `json:"address_line2"`
		City                  string `json:"city"`
		State                 string `json:"state"`
		PostalCode            string `json:"postal_code"`
		Country               string `json:"country"`
		EmergencyContactName  string `json:"emergency_contact_name"`
		EmergencyContactPhone string `json:"emergency_contact_phone"`
		EmergencyContactRelation string `json:"emergency_contact_relation"`
		EmployeeNumber        string `json:"employee_number"`
		HireDate              *string `json:"hire_date"`
		TerminationDate       *string `json:"termination_date"`
		Status                string `json:"status"`
		DepartmentID          *uint  `json:"department_id"`
		PositionID            *uint  `json:"position_id"`
		ManagerID             *uint  `json:"manager_id"`
		BankAccountNumber     string `json:"bank_account_number"`
		BankName              string `json:"bank_name"`
		BankRoutingNumber     string `json:"bank_routing_number"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	employee := &model.Employee{
		FirstName:             body.FirstName,
		LastName:              body.LastName,
		Email:                 body.Email,
		Phone:                 body.Phone,
		Gender:                body.Gender,
		AddressLine1:          body.AddressLine1,
		AddressLine2:          body.AddressLine2,
		City:                  body.City,
		State:                 body.State,
		PostalCode:            body.PostalCode,
		Country:               body.Country,
		EmergencyContactName:  body.EmergencyContactName,
		EmergencyContactPhone: body.EmergencyContactPhone,
		EmergencyContactRelation: body.EmergencyContactRelation,
		EmployeeNumber:        body.EmployeeNumber,
		Status:                body.Status,
		DepartmentID:          body.DepartmentID,
		PositionID:            body.PositionID,
		ManagerID:             body.ManagerID,
		BankAccountNumber:     body.BankAccountNumber,
		BankName:              body.BankName,
		BankRoutingNumber:     body.BankRoutingNumber,
	}

	if err := h.service.CreateEmployee(1, employee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Employee created successfully"})
}

// GetEmployeeByID handles GET /employees/:id
// @Summary Get employee by ID
// @Description Get detailed information about an employee
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Security ApiKeyAuth
// @Success 200 {object} model.Employee
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/employees/{id} [get]
func (h *EmployeeHandler) GetEmployeeByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	employee, err := h.service.GetEmployeeByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employee not found"})
	}
	return c.Status(fiber.StatusOK).JSON(employee)
}

func (h *EmployeeHandler) GetEmployeeByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	employee, err := h.service.GetEmployeeByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employee not found"})
	}
	return c.Status(fiber.StatusOK).JSON(employee)
}

// GetAllEmployees handles GET /employees
// @Summary Get employees with filtering
// @Description Get a list of employees with pagination and optional filters (email, status, department_id, manager_id, search)
// @Tags Employees
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Param email query string false "Filter by email"
// @Param status query string false "Filter by status"
// @Param department_id query int false "Filter by department ID"
// @Param manager_id query int false "Filter by manager ID"
// @Param search query string false "Search by name"
// @Security ApiKeyAuth
// @Success 200 {array} model.Employee
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *fiber.Ctx) error {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	email := c.Query("email")
	status := c.Query("status")
	deptIDStr := c.Query("department_id")
	mgrIDStr := c.Query("manager_id")
	search := c.Query("search")

	limit, _ := strconv.Atoi(limitStr)
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	filter := repository.EmployeeFilter{
		Email:  email,
		Status: status,
		Search: search,
	}

	if deptIDStr != "" {
		id, err := strconv.ParseUint(deptIDStr, 10, 32)
		if err == nil {
			uID := uint(id)
			filter.DepartmentID = &uID
		}
	}
	if mgrIDStr != "" {
		id, err := strconv.ParseUint(mgrIDStr, 10, 32)
		if err == nil {
			uID := uint(id)
			filter.ManagerID = &uID
		}
	}

	employees, err := h.service.GetAllEmployees(limit, offset, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(employees)
}

// UpdateEmployee handles PUT /employees/:id
// @Summary Update an employee
// @Description Update an existing employee's details
// @Tags Employees
// @Accept json
// @Produce json
// @Param id path int true "Employee ID"
// @Param employee body model.Employee true "Employee update object"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}

	var employee model.Employee
	if err := c.BodyParser(&employee); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	employee.ID = uint(id)

	if err := h.service.UpdateEmployee(1, &employee); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Employee updated successfully"})
}

// DeleteEmployee handles DELETE /employees/:id
// @Summary Delete an employee
// @Description Delete an employee by their ID
// @Tags Employees
// @Produce json
// @Param id path int true "Employee ID"
// @Security ApiKeyAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee ID"})
	}
	if err := h.service.DeleteEmployee(1, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Employee deleted successfully"})
}
