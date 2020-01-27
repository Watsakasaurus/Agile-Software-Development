package controllers

import (
	"medicare-api/db"
	"medicare-api/types"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

func GetMedicalData(c echo.Context) error {
	log.Debugf("GetMedicareData")
	dbc := c.Get("db").(*db.Client)
	perPage := c.Get("per_page").(int)
	pageNumber := c.Get("page_number").(int)

	filter := db.MedicalDataFilter{
		PriceMin:  getOptionalInt(c, "max_price"),
		PriceMax:  getOptionalInt(c, "min_price"),
		Proximity: getOptionalInt(c, "proximity"),
		Query:     getOptionalString(c, "query"),
		Latitude:  getOptionalFloat64(c, "lat"),
		Longitude: getOptionalFloat64(c, "long"),
	}

	// If proximity is not set, default to 200 miles
	if filter.Proximity == nil {
		proximity := 200
		filter.Proximity = &proximity
	}

	var results []types.MedicalData
	var dbErr *types.Error
	total := 0
	// If Latitude and Longitude is supplied query by location
	if filter.Latitude != nil && filter.Longitude != nil {
		results, total, dbErr = dbc.GetMedicalDataByLocation(filter, perPage, pageNumber)
	} else {
		// Otherwise search by description
		results, total, dbErr = dbc.GetMedicalDataByDescription(filter, perPage, pageNumber)
	}
	if dbErr != nil {
		return dbErr
	}

	return c.JSON(http.StatusOK, types.MedicareDataResponse{
		Objects:    results,
		Total:      total,
		PerPage:    perPage,
		PageNumber: pageNumber,
	})
}
