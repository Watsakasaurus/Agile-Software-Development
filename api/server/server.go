package server

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	sv "github.com/Rekfuki/swag-validator"
	"github.com/fvbock/endless"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/miketonks/swag"
	"github.com/miketonks/swag/swagger"
	log "github.com/sirupsen/logrus"

	"medicare-api/controllers"
	"medicare-api/db"
	"medicare-api/types"
)

// ContextParams stores context parameters for server initialization
type ContextParams struct {
	DB *db.Client
}

// ContextObjects attaches backend clients to the API context
func ContextObjects(contextParams ContextParams) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", contextParams.DB)
			return next(c)
		}
	}
}

// DefaultContentType sets default content type
func DefaultContentType() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			t := c.Request().Header.Get("Content-Type")
			if t == "" {
				c.Request().Header.Set("Content-Type", "application/json; charset=UTF-8")
			}
			return next(c)
		}
	}
}

// CreateRouter creates the router.
func CreateRouter(params ContextParams) *echo.Echo {
	r := echo.New()
	r.Debug = true
	r.HTTPErrorHandler = errorHandler

	r.Use(
		ContextObjects(params),
		middleware.Logger(),
		middleware.Recover(),
		middleware.Secure(),
		middleware.CORS(),
	)

	r.GET("/ping", controllers.Ping)

	medicareAPI := CreateSwaggerAPI()

	// Swagger UI
	r.GET("/medicare/api/json", echo.WrapHandler(medicareAPI.Handler(true)))

	api := r.Group("", sv.SwaggerValidatorEcho(medicareAPI), DefaultContentType(), Pagination())
	medicareAPI.Walk(func(path string, endpoint *swagger.Endpoint) {
		h := endpoint.Handler.(func(c echo.Context) error)
		path = swag.ColonPath(path)
		api.Add(endpoint.Method, path, h)
	})

	var routes []string
	for _, route := range r.Routes() {
		if route.Path != "." && route.Path != "/*" {
			routes = append(routes, route.Method+" "+route.Path)
		}
	}
	sort.Strings(routes)
	for _, route := range routes {
		log.Print(route)
	}

	return r
}

// CreateSwaggerAPI creates all swagger endpoints.
func CreateSwaggerAPI() *swagger.API {
	api := swag.New(
		swag.Title("Medicare API"),
		swag.Version("2.0"),
		swag.BasePath("/medicare/api"),
		swag.Endpoints(
			aggregateEndpoints()...,
		),
	)
	return api
}

func aggregateEndpoints(endpoints ...[]*swagger.Endpoint) []*swagger.Endpoint {
	res := []*swagger.Endpoint{}
	for _, v := range endpoints {
		res = append(res, v...)
	}
	return res
}

// Run runs the server
func Run(params ContextParams) {
	r := CreateRouter(params)

	endless.DefaultHammerTime = 10 * time.Second
	endless.DefaultReadTimeOut = 295 * time.Second
	if err := endless.ListenAndServe(":5009", r); err != nil {
		log.Infof("Server stopped: %s", err)
	}
}

func errorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := err.Error()
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	} else if te, ok := err.(*types.Error); ok {
		switch te.Type {
		case types.ErrTypeValidationError:
			code = http.StatusBadRequest
		case types.ErrTypeDuplicateError:
			code = http.StatusConflict
		case types.ErrTypeNotFoundError:
			code = http.StatusNotFound
		}
	}
	log.Error(fmt.Sprintf("%d %s", code, message))
	err2 := c.JSON(code, types.ErrorResponse{Message: message})
	if err2 != nil {
		log.Errorf("failed to handle error: %s", err2)
	}
}
