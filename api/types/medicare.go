package types

import (
	"encoding/json"
	"fmt"
)

// MedicalData represents medical data
type MedicalData struct {
	ProviderID           int     `db:"provider_id" json:"provider_id"`
	ProviderName         string  `db:"provider_name" json:"provider_name"`
	ProviderStreet       string  `db:"provider_street" json:"provider_street"`
	ProviderCity         string  `db:"provider_city" json:"provider_city"`
	ProviderState        string  `db:"provider_state" json:"provider_state"`
	ProviderZipCode      int     `db:"provider_zip_code" json:"provider_zip_code"`
	HRRDescription       string  `db:"hrr_description" json:"hrr_description"`
	TotalDischarges      int     `db:"total_discharges" json:"total_discharges"`
	AverageTotalPayments float64 `db:"average_total_payments" json:"average_total_payments"`
	DRGDefinition        string  `db:"drg_definition" json:"drg_definition"`
	Distance             float64 `db:"distance" json:"distance,omitempty"`
	Latitude             float64 `db:"latitude" json:"latitude,omitempty"`
	Longitude            float64 `db:"longitude" json:"longitude,omitempty"`
}

// MedicareDataResponse represents response for GET ALL
type MedicareDataResponse struct {
	Objects    []MedicalData `json:"objects"`
	Total      int           `json:"total"`
	PageNumber int           `json:"page_number"`
	PerPage    int           `json:"per_page"`
}

type DrgDefinitions []string

// ProcedureForFiltering is minimal amount of data necessary to filter
type FilteringData struct {
	DrgDefinitions DrgDefinitions `db:"drg_definitions" json:"procedure_definitions"`
	PriceMin       float64        `db:"price_min" json:"price_min"`
	PriceMax       float64        `db:"price_max" json:"price_max"`
}

func (d *DrgDefinitions) Scan(value interface{}) error {
	switch src := value.(type) {
	case []byte:
		return json.Unmarshal(src, d)
	case nil:
		*d = []string{}
		return nil
	default:
		return fmt.Errorf("invalid type for Document: %T", src)
	}
}
