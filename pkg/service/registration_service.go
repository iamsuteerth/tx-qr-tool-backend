package service

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/dto"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/models"
	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/repository"
	"github.com/iamsuteerth/tx-qr-tool-backend/utils"
)

type RegistrationService interface {
	CreateRegistration(req *dto.CreateRegistrationRequest) (*dto.RegistrationResponse, error)
	GenerateCSV() ([]byte, error)
}

type registrationService struct {
	repo repository.RegistrationRepository
}

func NewRegistrationService(repo repository.RegistrationRepository) RegistrationService {
	return &registrationService{repo: repo}
}

func (s *registrationService) CreateRegistration(req *dto.CreateRegistrationRequest) (*dto.RegistrationResponse, error) {
	// Validate request
	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		return nil, &utils.AppError{
			HTTPCode:         400,
			Code:             "VALIDATION_ERROR",
			Message:          "Invalid input data",
			ValidationErrors: validationErrors,
		}
	}

	// Check for duplicate email
	existingReg, err := s.repo.GetByEmail(req.Email)
	if err == nil && existingReg != nil {
		return nil, utils.NewBadRequestError("DUPLICATE_EMAIL", "Email already registered", nil)
	}

	// Check for duplicate phone
	existingRegByPhone, err := s.repo.GetByPhone(req.Phone) // Phone is now string
	if err == nil && existingRegByPhone != nil {
		return nil, utils.NewBadRequestError("DUPLICATE_PHONE", "Phone number already registered", nil)
	}

	// Create registration model
	registration := &models.Registration{
		FullName:    req.FullName,
		Email:       req.Email,
		Phone:       req.Phone, // Phone is now string
		OrgName:     req.OrgName,
		Designation: req.Designation,
		MktSource:   req.MktSource,
		FoodPref:    req.FoodPref,
		TShirt:      req.TShirt,
	}

	// Save to database
	createdReg, err := s.repo.Create(registration)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	return &dto.RegistrationResponse{
		ID:          createdReg.ID,
		FullName:    createdReg.FullName,
		Email:       createdReg.Email,
		Phone:       createdReg.Phone, // Phone is now string
		OrgName:     createdReg.OrgName,
		Designation: createdReg.Designation,
		MktSource:   createdReg.MktSource,
		FoodPref:    createdReg.FoodPref,
		TShirt:      createdReg.TShirt,
		CreatedOn:   createdReg.CreatedOn.Format(time.RFC3339),
	}, nil
}

func (s *registrationService) GenerateCSV() ([]byte, error) {
	registrations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header
	header := []string{
		"ID",
		"Full Name",
		"Email",
		"Phone",
		"Organization",
		"Designation",
		"Marketing Source",
		"Food Preference",
		"T-Shirt Size",
		"Created On",
	}
	if err := writer.Write(header); err != nil {
		return nil, utils.NewInternalServerError("CSV_ERROR", "Failed to write CSV header", err)
	}

	// Write data
	for _, reg := range registrations {
		record := []string{
			strconv.Itoa(reg.ID),
			reg.FullName,
			reg.Email,
			reg.Phone, // Phone is now string, no conversion needed
			reg.OrgName,
			reg.Designation,
			reg.MktSource,
			reg.FoodPref,
			reg.TShirt,
			reg.CreatedOn.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(record); err != nil {
			return nil, utils.NewInternalServerError("CSV_ERROR", "Failed to write CSV record", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, utils.NewInternalServerError("CSV_ERROR", "Failed to generate CSV", err)
	}

	return buf.Bytes(), nil
}
