package dto

import (
	"strings"

	"github.com/iamsuteerth/tx-qr-tool-backend/utils"
)

type CreateRegistrationRequest struct {
	FullName    string `json:"full_name" binding:"required,min=2,max=255"`
	Email       string `json:"email" binding:"required"`
	Phone       string `json:"phone" binding:"required,min=10,max=13"`
	OrgName     string `json:"org_name,omitempty"`
	Designation string `json:"designation,omitempty"`
	MktSource   string `json:"mkt_source,omitempty"`
	FoodPref    string `json:"food_pref" binding:"required"`
	TShirt      string `json:"t-shirt" binding:"required,oneof=S M L XL XXL XXXL"`
}

func (r *CreateRegistrationRequest) Validate() []utils.ValidationError {
	var errors []utils.ValidationError

	if strings.TrimSpace(r.Email) == "" {
		errors = append(errors, utils.ValidationError{
			Field:   "email",
			Message: "Email cannot be empty",
		})
	}

	if strings.TrimSpace(r.FullName) == "" {
		errors = append(errors, utils.ValidationError{
			Field:   "full_name",
			Message: "Name cannot be empty",
		})
	}

	return errors
}

type RegistrationResponse struct {
	ID          int    `json:"id"`
	FullName    string `json:"full_name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	OrgName     string `json:"org_name"`
	Designation string `json:"designation"`
	MktSource   string `json:"mkt_source"`
	FoodPref    string `json:"food_pref"`
	TShirt      string `json:"t_shirt"`
	CreatedOn   string `json:"created_on"`
}

type RegistrationListResponse struct {
	Registrations []RegistrationResponse `json:"registrations"`
	Total         int                    `json:"total"`
}
