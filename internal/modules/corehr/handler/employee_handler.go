package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/thecemakin/hr-project/internal/modules/corehr/model"
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

	if err := h.service.CreateEmployee(employee); err != nil {
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

// GetEmployeeByEmail handles GET /employees/email/:email
// @Summary Get employee by email
// @Description Get detailed information about an employee by their email address
// @Tags Employees
// @Produce json
// @Param email path string true "Employee Email"
// @Security ApiKeyAuth
// @Success 200 {object} model.Employee
// @Failure 404 {object} map[string]string
// @Router /api/v1/corehr/employees/email/{email} [get]
func (h *EmployeeHandler) GetEmployeeByEmail(c *fiber.Ctx) error {
	email := c.Params("email")
	employee, err := h.service.GetEmployeeByEmail(email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Employee not found"})
	}
	return c.Status(fiber.StatusOK).JSON(employee)
}

// GetAllEmployees handles GET /employees
// @Summary Get all employees
// @Description Get a list of all employees with pagination
// @Tags Employees
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Security ApiKeyAuth
// @Success 200 {array} model.Employee
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees [get]
func (h *EmployeeHandler) GetAllEmployees(c *fiber.Ctx) error {
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

	employees, err := h.service.GetAllEmployees(limit, offset)
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

	if err := h.service.UpdateEmployee(&employee); err != nil {
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
	if err := h.service.DeleteEmployee(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Employee deleted successfully"})
}

// GetEmployeesByManagerID handles GET /employees/manager/:managerId
// @Summary Get employees by manager ID
// @Description Get a list of employees reporting to a specific manager
// @Tags Employees
// @Produce json
// @Param managerId path int true "Manager ID"
// @Security ApiKeyAuth
// @Success 200 {array} model.Employee
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees/manager/{managerId} [get]
func (h *EmployeeHandler) GetEmployeesByManagerID(c *fiber.Ctx) error {
	idParam := c.Params("managerId")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid manager ID"})
	}
	employees, err := h.service.GetEmployeesByManagerID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(employees)
}

// GetEmployeesByDepartmentID handles GET /employees/department/:departmentId
// @Summary Get employees by department ID
// @Description Get a list of employees in a specific department
// @Tags Employees
// @Produce json
// @Param departmentId path int true "Department ID"
// @Security ApiKeyAuth
// @Success 200 {array} model.Employee
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees/department/{departmentId} [get]
func (h *EmployeeHandler) GetEmployeesByDepartmentID(c *fiber.Ctx) error {
	idParam := c.Params("departmentId")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid department ID"})
	}
	employees, err := h.service.GetEmployeesByDepartmentID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(employees)
}

// GetEmployeesByStatus handles GET /employees/status/:status
// @Summary Get employees by status
// @Description Get a list of employees with a specific employment status
// @Tags Employees
// @Produce json
// @Param status path string true "Status (active, inactive, terminated)"
// @Security ApiKeyAuth
// @Success 200 {array} model.Employee
// @Failure 500 {object} map[string]string
// @Router /api/v1/corehr/employees/status/{status} [get]
func (h *EmployeeHandler) GetEmployeesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	employees, err := h.service.GetEmployeesByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(employees)
}