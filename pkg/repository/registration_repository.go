package repository

import (
	"context"
	"errors"

	"github.com/iamsuteerth/tx-qr-tool-backend/pkg/models"
	"github.com/iamsuteerth/tx-qr-tool-backend/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RegistrationRepository interface {
	Create(registration *models.Registration) (*models.Registration, error)
	GetAll() ([]models.Registration, error)
	GetByEmail(email string) (*models.Registration, error)
	GetByPhone(phone string) (*models.Registration, error) // Changed parameter type
}

type registrationRepository struct {
	db *pgxpool.Pool
}

func NewRegistrationRepository(db *pgxpool.Pool) RegistrationRepository {
	return &registrationRepository{db: db}
}

func (r *registrationRepository) Create(registration *models.Registration) (*models.Registration, error) {
	query := `
        INSERT INTO registrations (full_name, email, phone, org_name, designation, mkt_source, food_pref, t_shirt)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_on
    `

	ctx := context.Background()
	err := r.db.QueryRow(
		ctx,
		query,
		registration.FullName,
		registration.Email,
		registration.Phone,
		registration.OrgName,
		registration.Designation,
		registration.MktSource,
		registration.FoodPref,
		registration.TShirt,
	).Scan(&registration.ID, &registration.CreatedOn)

	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to create registration", err)
	}

	return registration, nil
}

func (r *registrationRepository) GetAll() ([]models.Registration, error) {
	query := `
        SELECT id, full_name, email, phone, org_name, designation, mkt_source, food_pref, t_shirt, created_on
        FROM registrations
        ORDER BY created_on DESC
    `

	ctx := context.Background()
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to fetch registrations", err)
	}
	defer rows.Close()

	var registrations []models.Registration
	for rows.Next() {
		var reg models.Registration
		err := rows.Scan(
			&reg.ID,
			&reg.FullName,
			&reg.Email,
			&reg.Phone,
			&reg.OrgName,
			&reg.Designation,
			&reg.MktSource,
			&reg.FoodPref,
			&reg.TShirt,
			&reg.CreatedOn,
		)
		if err != nil {
			return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to scan registration", err)
		}
		registrations = append(registrations, reg)
	}

	if err = rows.Err(); err != nil {
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Error iterating registrations", err)
	}

	return registrations, nil
}

func (r *registrationRepository) GetByEmail(email string) (*models.Registration, error) {
	query := `
        SELECT id, full_name, email, phone, org_name, designation, mkt_source, food_pref, t_shirt, created_on
        FROM registrations
        WHERE email = $1
    `

	ctx := context.Background()
	var reg models.Registration
	err := r.db.QueryRow(ctx, query, email).Scan(
		&reg.ID,
		&reg.FullName,
		&reg.Email,
		&reg.Phone,
		&reg.OrgName,
		&reg.Designation,
		&reg.MktSource,
		&reg.FoodPref,
		&reg.TShirt,
		&reg.CreatedOn,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewNotFoundError("REGISTRATION_NOT_FOUND", "Registration not found", err)
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to fetch registration", err)
	}

	return &reg, nil
}

func (r *registrationRepository) GetByPhone(phone string) (*models.Registration, error) { 
	query := `
        SELECT id, full_name, email, phone, org_name, designation, mkt_source, food_pref, t_shirt, created_on
        FROM registrations
        WHERE phone = $1
    `

	ctx := context.Background()
	var reg models.Registration
	err := r.db.QueryRow(ctx, query, phone).Scan(
		&reg.ID,
		&reg.FullName,
		&reg.Email,
		&reg.Phone,
		&reg.OrgName,
		&reg.Designation,
		&reg.MktSource,
		&reg.FoodPref,
		&reg.TShirt,
		&reg.CreatedOn,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewNotFoundError("REGISTRATION_NOT_FOUND", "Registration not found", err)
		}
		return nil, utils.NewInternalServerError("DATABASE_ERROR", "Failed to fetch registration", err)
	}

	return &reg, nil
}
