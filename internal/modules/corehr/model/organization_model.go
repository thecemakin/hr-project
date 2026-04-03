package model

// OrganizationNode represents a node in the organizational hierarchy tree
type OrganizationNode struct {
	ID           uint               `json:"id"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	JobTitle     string             `json:"job_title"`
	Department   string             `json:"department"`
	Status       string             `json:"status"`
	Subordinates []*OrganizationNode `json:"subordinates,omitempty"`
}
