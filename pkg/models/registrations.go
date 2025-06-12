package models

import (
	"time"
)

type Registration struct {
	ID          int       `json:"id" db:"id"`
	FullName    string    `json:"full_name" db:"full_name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	OrgName     string    `json:"org_name" db:"org_name"`
	Designation string    `json:"designation" db:"designation"`
	MktSource   string    `json:"mkt_source" db:"mkt_source"`
	FoodPref    string    `json:"food_pref" db:"food_pref"`
	TShirt      string    `json:"t_shirt" db:"t_shirt"`
	CreatedOn   time.Time `json:"created_on" db:"created_on"`
}
