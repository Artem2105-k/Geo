package controller

import (
	"net/http"
	"strings"

	"ar.konovalov202_gmail.com/rpc/general"
	"ar.konovalov202_gmail.com/rpc/internal/models"
	"github.com/go-chi/render"
)

type GeoController struct {
	service general.GeoProvider
}

func NewGeoController(s general.GeoProvider) *GeoController {
	return &GeoController{service: s}
}

// SearchRequest определяет структуру запроса
type SearchRequest struct {
	Query string `json:"query"`
}

// GeocodeRequest определяет структуру запроса геокодирования
type GeocodeRequest struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// AddressResponse определяет структуру ответа
type AddressResponse struct {
	Addresses []*general.Address `json:"addresses"`
}

// SearchAddress godoc
// @Summary Поиск адреса
// @Description Поиск адресов по строке запроса
// @Tags address
// @Accept json
// @Produce json
// @Param input body SearchRequest true "Поисковый запрос"
// @Success 200 {object} AddressResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/address/search [post]
func (c *GeoController) SearchAddress(w http.ResponseWriter, r *http.Request) {
	var req SearchRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		models.RenderError(w, r, "Invalid request format", http.StatusBadRequest)
		return
	}

	req.Query = strings.TrimSpace(req.Query)
	if req.Query == "" {
		models.RenderError(w, r, "Query parameter is required", http.StatusBadRequest)
		return
	}

	addresses, err := c.service.AddressSearch(req.Query)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, AddressResponse{Addresses: addresses})
}

// GeoCode godoc
// @Summary Геокодирование
// @Description Поиск адреса по координатам
// @Tags address
// @Accept json
// @Produce json
// @Param input body GeocodeRequest true "Координаты"
// @Success 200 {object} AddressResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/address/geocode [post]
func (c *GeoController) GeoCode(w http.ResponseWriter, r *http.Request) {
	var req GeocodeRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		models.RenderError(w, r, "Invalid request format", http.StatusBadRequest)
		return
	}

	if req.Lat == "" || req.Lng == "" {
		models.RenderError(w, r, "Lat and Lng parameters are required", http.StatusBadRequest)
		return
	}

	addresses, err := c.service.GeoCode(req.Lat, req.Lng)
	if err != nil {
		handleServiceError(w, r, err)
		return
	}

	render.Status(r, http.StatusOK)
	render.JSON(w, r, AddressResponse{Addresses: addresses})
}

func handleServiceError(w http.ResponseWriter, r *http.Request, err error) {
	if strings.Contains(err.Error(), "connection failed") ||
		strings.Contains(err.Error(), "status code: 5") {
		models.RenderError(w, r, "Dadata service unavailable", http.StatusInternalServerError)
	} else {
		models.RenderError(w, r, err.Error(), http.StatusBadRequest)
	}
}
