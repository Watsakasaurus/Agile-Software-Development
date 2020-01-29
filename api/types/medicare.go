package types

import (
	"encoding/json"
	"fmt"
)

// MedicalData represents medical data
type MedicalData struct {
	ProviderID              int     `db:"provider_id" json:"provider_id"`
	ProviderName            string  `db:"provider_name" json:"provider_name"`
	ProviderStreet          string  `db:"provider_street" json:"provider_street"`
	ProviderCity            string  `db:"provider_city" json:"provider_city"`
	ProviderState           string  `db:"provider_state" json:"provider_state"`
	ProviderZipCode         int     `db:"provider_zip_code" json:"provider_zip_code"`
	HRRDescription          string  `db:"hrr_description" json:"hrr_description"`
	TotalDischarges         int     `db:"total_discharges" json:"-"`
	AverageCoveredCharges   float64 `db:"average_covered_charges" json:"average_covered_charges"`
	AverageTotalPayments    float64 `db:"average_total_payments" json:"average_total_payments"`
	AverageMedicarePayments float64 `db:"average_medicare_payments" json:"average_medicare_payments"`
	DRGDefinition           string  `db:"drg_definition" json:"drg_definition"`
	DRGDefinitionTokens     string  `db:"drg_definition_tokens" json:"-"`
	Distance                float64 `db:"distance" json:"distance,omitempty"`
	Latitude                float64 `db:"latitude" json:"latitude,omitempty"`
	Longitude               float64 `db:"longitude" json:"longitude,omitempty"`
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

// ZipCodeLatLong represents zip_code_lat_long table row
type ZipCodeLatLong struct {
	ZipCode   int     `db:"zip_code"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitudes"`
}

// Provider represents provider
type Provider struct {
	ID                int    `db:"id" json:"provider_id"`
	Name              string `db:"name" json:"provider_name"`
	Street            string `db:"street" json:"provider_street"`
	City              string `db:"city" json:"provider_city"`
	State             string `db:"state" json:"provider_state"`
	ZipCode           int    `db:"zip_code" json:"provider_zip_code"`
	HRRDescription    string `db:"hrr_description" json:"hrr_description"`
	AdressLine        string `db:"address_line" json:"-"`
	AddressLineTokens string `db:"address_line_tokens" json:"-"`
}

// Procedure represents procedure
type Procedure struct {
	ID                      string  `db:"id" json:"id"`
	TotalDischarges         int     `db:"total_discharges" json:"-"`
	AverageCoveredCharges   float64 `db:"average_covered_charges" json:"average_covered_charges"`
	AverageTotalPayments    float64 `db:"average_total_payments" json:"average_total_payments"`
	AverageMedicarePayments float64 `db:"average_medicare_payments" json:"average_medicare_payments"`
	DRGDefinition           string  `db:"drg_definition" json:"drg_definition"`
	DRGDefinitionTokens     string  `db:"drg_definition_tokens" json:"-"`
}

// ProviderProcedures represents a link between a provider and procedure
type ProvideRrocedure struct {
	ProviderID           int     `db:"provider_id" json:"provider_id"`
	ProcedureID          string  `db:"procedure_id" json:"procedure_id"`
	AverageTotalPayments float64 `db:"average_total_payments" json:"average_total_payments"`
}
