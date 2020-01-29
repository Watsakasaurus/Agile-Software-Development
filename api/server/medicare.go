package server

import (
	"net/http"

	"github.com/miketonks/swag/endpoint"
	"github.com/miketonks/swag/swagger"

	"medicare-api/controllers"
	"medicare-api/types"
)

func medicareAPI() []*swagger.Endpoint {
	getMedicalData := endpoint.New("GET", "/data", "Get medical data",
		endpoint.Handler(controllers.GetMedicalData),
		endpoint.Description("Returns medical data"),
		endpoint.Response(http.StatusOK, types.MedicareDataResponse{}, "Success"),
		endpoint.QueryMap(map[string]swagger.Parameter{
			"page": {
				Type:        "integer",
				Description: "Page number to return",
			},
			"per_page": {
				Type:        "integer",
				Description: "Number of records per page",
			},
			"min_price": {
				Type:        "integer",
				Minimum:     &[]int64{0}[0],
				Description: "Minimum price to search by",
			},
			"max_price": {
				Type:        "integer",
				Minimum:     &[]int64{1}[0],
				Description: "Maximum price to search by",
			},
			"proximity": {
				Type:        "integer",
				Minimum:     &[]int64{1}[0],
				Description: "Distance radius by which to limit return results",
			},
			"query": {
				Type:        "string",
				Description: "Text to query by",
			},
			"lat": {
				Type:        "number",
				Description: "Latitude",
			},
			"long": {
				Type:        "number",
				Description: "longitude",
			},
		}),
		endpoint.Tags("Medicare"),
	)
	getFilteringData := endpoint.New("GET", "/filtering", "Get data to perform filtering",
		endpoint.Handler(controllers.GetFilteringData),
		endpoint.Description("Returns data for filtering"),
		endpoint.Response(http.StatusOK, types.FilteringData{}, "Success"),
		endpoint.Tags("Medicare"),
	)

	return []*swagger.Endpoint{
		getMedicalData,
		getFilteringData,
	}
}
