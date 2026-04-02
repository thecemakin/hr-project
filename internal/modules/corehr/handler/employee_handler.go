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

func (h *EmployeeHandler) GetEmployeesByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	employees, err := h.service.GetEmployeesByStatus(status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(employees)
}