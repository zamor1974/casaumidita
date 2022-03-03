package controllers

import (
	"casaumidita/lang"
	"casaumidita/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
)

// BaseHandler will hold everything that controller needs
type BaseHandlerSqlx struct {
	db *sqlx.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandlerSqlx(db *sqlx.DB) *BaseHandlerSqlx {
	return &BaseHandlerSqlx{
		db: db,
	}
}

// swagger:model CommonError
type CommonError struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model CommonSuccess
type CommonSuccess struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:model GetHumidities
type GetHumidities struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string             `json:"message"`
	Data    *models.Humidities `json:"data"`
}

// swagger:model GetHumidity
type GetHumidity struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the response
	// in: string
	Message string `json:"message"`
	// Umidity value
	Data *models.Humidity `json:"data"`
}

// ErrHandler returns error message response
func ErrHandler(errmessage string) *CommonError {
	errresponse := CommonError{}
	errresponse.Status = 0
	errresponse.Message = errmessage
	return &errresponse
}

// swagger:route GET /humidities listhumidity
// Get humidity list
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetHumidities
func (h *BaseHandlerSqlx) GetHumiditiesSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetHumidities{}

	humidities := models.GetHumiditiesSqlx(h.db.DB)

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = humidities

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route GET /lasthour lastHour
// Get list of last hour of humidity values .... or the last value inserted
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetHumidities
func (h *BaseHandlerSqlx) GetLastHourSqlx(w http.ResponseWriter, r *http.Request) {
	response := GetHumidities{}

	companies := models.GetLastHourSqlx(h.db.DB)

	response.Status = 1
	response.Message = lang.Get("success")
	response.Data = companies

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// swagger:route POST /humidity addHumidity
// Create a new humidity value
//
// security:
// - apiKey: []
// responses:
//  401: CommonError
//  200: GetHumidity
func (h *BaseHandlerSqlx) PostHumiditySqlx(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	response := GetHumidity{}

	decoder := json.NewDecoder(r.Body)
	var reqHumidity *models.ReqAddHumidity
	err := decoder.Decode(&reqHumidity)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(ErrHandler(lang.Get("invalid_request")))
		return
	}

	rain, errmessage := models.PostHumiditySqlx(h.db.DB, reqHumidity)
	if errmessage != "" {
		json.NewEncoder(w).Encode(ErrHandler(errmessage))
		return
	}

	response.Status = 1
	response.Message = lang.Get("insert_success")
	response.Data = rain
	json.NewEncoder(w).Encode(response)
}
