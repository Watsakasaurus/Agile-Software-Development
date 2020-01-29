package db

import (
	"database/sql"
	"medicare-api/types"
	"strings"

	log "github.com/sirupsen/logrus"
	"xorm.io/builder"
)

type MedicalDataFilter struct {
	PriceMax  *int
	PriceMin  *int
	Proximity *int
	Query     *string
	Latitude  *float64
	Longitude *float64
}

func (c *Client) GetMedicalDataByLocation(filter MedicalDataFilter, perPage, pageNumber int) ([]types.MedicalData, int, *types.Error) {
	log.Debugf("GetMedicalDataByLocation")

	query := `
		select * from ( select
			p.id as provider_id, p.name as provider_name, p.street as provider_street, 
			p.city as provider_city, p.state as provider_state, p.zip_code as provider_zip_code,
			p.hrr_description, pr.drg_definition, pp.average_total_payments,
			pr.total_discharges, z.latitude, z.longitude,
			round(
				point(z.latitude, z.longitude)<@>point(?,?)
			) as distance
		from procedures pr
		join provider_procedures pp on pp.procedure_id=pr.id
		join providers p on p.id=pp.provider_id
		join zip_code_lat_long z on z.zip_code=p.zip_code
		order by distance) as res where res.distance<=?	
	`

	args := []interface{}{*filter.Latitude, *filter.Longitude, *filter.Proximity}
	if filter.PriceMax != nil {
		query += " and res.average_total_payments <= ?"
		args = append(args, *filter.PriceMax)
	}
	if filter.PriceMin != nil {
		query += " and res.average_total_payments >= ?"
		args = append(args, *filter.PriceMin)
	}
	if filter.Query != nil {
		tokens := strings.Split(*filter.Query, " ")
		joined := strings.Join(tokens, "%")
		joined = "%" + joined + "%"

		query += " and res.drg_definition ilike ?"
		args = append(args, joined)
	}

	orderBy := "res.distance ASC"

	query = c.ex.Rebind(query)
	results := []types.MedicalData{}
	total, err := c.SelectWithCountSQL(&results, query, args, orderBy, pageNumber, perPage)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, c.transformError(err)
	}

	return results, total, nil
}

func (c *Client) GetMedicalDataByDescription(filter MedicalDataFilter, perPage, pageNumber int) ([]types.MedicalData, int, *types.Error) {
	log.Debugf("GetMedicalDataByDescription")

	query := c.Builder().
		Select(`
		p.id as provider_id, p.name as provider_name, p.street as provider_street, 
		p.city as provider_city, p.state as provider_state, p.zip_code as provider_zip_code,
		p.hrr_description, pr.drg_definition, pr.total_discharges, pp.average_total_payments,
		zcll.latitude, zcll.longitude`).
		From("procedures pr").
		InnerJoin("provider_procedures pp", "pp.procedure_id=pr.id").
		InnerJoin("providers p", "p.id=pp.provider_id").
		InnerJoin("zip_code_lat_long zcll", "zcll.zip_code=p.zip_code").
		OrderBy("pp.average_total_payments ASC, p.name ASC")

	// nolint: unparam
	query = applyMedicalDataFilter(query, filter)

	results := []types.MedicalData{}

	total, err := c.SelectWithCount(&results, query, pageNumber, perPage)
	if err != nil && err != sql.ErrNoRows {
		return nil, 0, c.transformError(err)
	}

	return results, total, nil
}

// GetProceduresForFiltering returns bare minimum information for filtering
func (c *Client) GetFilteringData() (*types.FilteringData, *types.Error) {
	log.Debugf("GetFilteringData")

	query := c.Builder().
		Select(`to_json(array(select distinct pr.drg_definition from provider_procedures pp join procedures pr on pr.id=pp.procedure_id))
		as drg_definitions, max(pp.average_total_payments) as price_max, min(pp.average_total_payments) as price_min`).
		From("provider_procedures pp")

	var results types.FilteringData
	err := c.Get(&results, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, c.transformError(err)
	}

	return &results, nil
}

// CreateProvider creates provider
func (c *Client) CreateProvider(payload types.Provider) *types.Error {
	log.Debugf("CreateProvider")

	query := `INSERT INTO providers 
		(id, name, street, city, state, zip_code, hrr_description)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := c.ex.Exec(query, payload.ID, payload.Name,
		payload.Street, payload.City, payload.State, payload.ZipCode,
		payload.HRRDescription)

	return c.transformError(err)
}

// func (c *Client) CreateProvider(payload types.Provider) (*types.Provider, *types.Error) {
// 	log.Debugf("CreateProvider")

// 	query := c.Builder().
// 		Insert(builder.Eq{
// 			"id":              payload.ID,
// 			"name":            payload.Name,
// 			"street":          payload.Street,
// 			"city":            payload.City,
// 			"state":           payload.State,
// 			"zip_code":        payload.ZipCode,
// 			"hrr_description": payload.HRRDescription,
// 		}).Into("providers")

// 	_, err := c.Exec(query)
// 	if err != nil {
// 		return nil, c.transformError(err)
// 	}

// 	return &payload, nil
// }

func (c *Client) GetProviderByID(id int) (*types.Provider, *types.Error) {
	log.Debugf("GetProviderByID")

	query := c.Builder().Select("*").
		From("providers").
		Where(builder.Eq{"id": id})

	var provider types.Provider
	err := c.Get(&provider, query)
	if err != nil {
		return nil, c.transformError(err)
	}

	return &provider, nil
}

// CreateProcedure inserts procedure into db
func (c *Client) CreateProcedure(payload types.Procedure) *types.Error {
	log.Debugf("CreateProcedure")

	// Insert procedure data
	query := `INSERT INTO procedures
		(total_discharges, drg_definition)
		VALUES ($1, $2)`

	_, err := c.ex.Exec(query, payload.TotalDischarges, payload.DRGDefinition)

	return c.transformError(err)
}

func (c *Client) GetProcedureByID(id string) (*types.Procedure, *types.Error) {
	log.Debugf("GetProcedureByID")

	query := c.Builder().Select("*").
		From("procedures").
		Where(builder.Eq{"id": id})

	var procedure types.Procedure
	err := c.Get(&procedure, query)
	if err != nil {
		return nil, c.transformError(err)
	}

	return &procedure, nil
}

// AssignProcedureToProvider creates a link between provider and procedure
func (c *Client) AssignProcedureToProvider(payload types.ProvideRrocedure) *types.Error {
	log.Debugf("AssignProcedureToProvider")

	query := `INSERT INTO provider_procedures (provider_id, procedure_id, average_total_payments) 
		VALUES($1, $2, $3)`

	_, err := c.ex.Exec(query, payload.ProviderID, payload.ProcedureID, payload.AverageTotalPayments)

	return c.transformError(err)
}

// CreateZipCodeLatLong assigns latitude and longitude to a zip code
func (c *Client) CreateZipCodeLatLong(payload types.ZipCodeLatLong) (*types.ZipCodeLatLong, *types.Error) {
	log.Debugf("CreateZipCodeLatLong")

	query := c.Builder().Insert(builder.Eq{
		"zip_code":  payload.ZipCode,
		"latitude":  payload.Latitude,
		"longitude": payload.Longitude,
	}).Into("zip_code_lat_long")

	_, err := c.Exec(query)
	if err != nil {
		return nil, c.transformError(err)
	}

	return &payload, nil
}

// apply filters to customer query, if required
// nolint: unparam
func applyMedicalDataFilter(query *builder.Builder, filter MedicalDataFilter) *builder.Builder {
	if filter.PriceMax != nil {
		query = query.And(builder.Lte{"pp.average_total_payments": *filter.PriceMax})
	}
	if filter.PriceMin != nil {
		query = query.And(builder.Gte{"pp.average_total_payments": *filter.PriceMin})
	}
	if filter.Query != nil {
		tokens := strings.Split(*filter.Query, " ")
		joined := strings.Join(tokens, "%")
		joined = "%" + joined + "%"

		query = query.And(builder.Expr("pr.drg_definition ilike ?", joined))
	}

	return query
}
