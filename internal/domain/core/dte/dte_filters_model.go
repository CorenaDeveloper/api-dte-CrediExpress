package dte

import "time"

type DTEFilters struct {
	BranchID     uint       `query:"-"`
	IncludeAll   bool       `query:"all,omitempty"`
	StartDate    *time.Time `query:"startDate,omitempty"`
	EndDate      *time.Time `query:"endDate,omitempty"`
	Status       string     `query:"status,omitempty"`
	Transmission string     `query:"transmission,omitempty"`
	DTEType      string     `query:"type,omitempty"`

	// Paginación
	Page     int `json:"page,omitempty"`
	PageSize int `json:"page_size,omitempty"`
}
