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
			p.hrr_description, pr.drg_definition, pr.drg_definition_tokens, pp.average_total_payments,
			pr.total_discharges, z.latitude, z.longitude,
			round(
				point(z.latitude, z.longitude)<@>point(?,?)
			) as distance
		from procedures pr
		join provider_procedures pp on pp.procedure_id=pr.id
		join providers p on p.id=pp.provider_id
		join zip_code_lat_long z on z.zip_code=p.zip_code
		order by distance) as res where res.distance<?	
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
		joined := strings.Join(tokens, " & ")

		query += " and res.drg_definition_tokens @@ ?::tsquery"
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
		OrderBy("average_total_payments ASC")

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
		Select(`
			to_json(array_agg(drg_definition)) as drg_definitions, MIN(CEIL(average_total_payments)) as price_min, 
			MAX(FLOOR(average_total_payments)) as price_max`).
		From("procedures")

	var results types.FilteringData
	err := c.Get(&results, query)
	if err != nil && err != sql.ErrNoRows {
		return nil, c.transformError(err)
	}

	return &results, nil
}

// apply filters to customer query, if required
// nolint: unparam
func applyMedicalDataFilter(query *builder.Builder, filter MedicalDataFilter) *builder.Builder {
	if filter.PriceMax != nil {
		query = query.And(builder.Gte{"pr.average_total_payments": *filter.PriceMax})
	}
	if filter.PriceMin != nil {
		query = query.And(builder.Lte{"pr.average_total_payments": *filter.PriceMin})
	}
	if filter.Query != nil {
		tokens := strings.Split(*filter.Query, " ")
		joined := strings.Join(tokens, " & ")

		query = query.And(builder.Expr("pr.drg_definition_tokens @@ ?::tsquery", joined))
	}

	return query
}
