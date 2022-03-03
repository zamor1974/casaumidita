package models

import (
	"casapioggia/lang"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// swagger:model Rain
type Rain struct {
	// Id of rain value
	// in: int64
	Id int64 `json:"id"`
	// Value of rain
	// in: int
	Value int `json:"valore"`
	// Timestamp of insert
	// in: time
	DateInsert time.Time `json:"data_inserimento"`
}

type Rains []Rain

type ReqAddRain struct {
	// Value of the rain
	// in: int
	Value int `json:"valore" validate:"required"`
}

// swagger:parameters addRain
type ReqRainBody struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddRain"
	//  required: true
	Body ReqAddRain `json:"body"`
}

// ErrHandler returns error message bassed on env debug
func ErrHandler(err error) string {
	var errmessage string
	if os.Getenv("DEBUG") == "true" {
		errmessage = err.Error()
	} else {
		errmessage = lang.Get("something_went_wrong")
	}
	return errmessage
}

func GetRainsSqlx(db *sql.DB) *Rains {
	rains := Rains{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM pioggia order by id desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Rain
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		rains = append(rains, p)
	}
	return &rains
}
func GetLastRainSqlx(db *sql.DB) *Rains {
	rains := Rains{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM pioggia where id = (select max(id) from pioggia)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Rain
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		rains = append(rains, p)
	}
	return &rains
}
func GetLastHourSqlx(db *sql.DB) *Rains {
	rains := Rains{}

	tFine := time.Now()
	dataFine := tFine.Format("2006-01-02 15:04:05")

	tInizio := time.Now().Add(time.Duration(-1) * time.Hour)
	dataInizio := tInizio.Format("2006-01-02 15:04:05")

	sqlStatement := fmt.Sprintf("SELECT id,valore,data_inserimento FROM pioggia where data_inserimento  >= '%s' AND data_inserimento <= '%s'", dataInizio, dataFine)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Rain
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		rains = append(rains, p)
	}

	if len(rains) == 0 {
		elemento := GetLastRainSqlx(db)
		rains = append(rains, *elemento...)
	}
	return &rains
}

// PostRainSqlx insert rain value
func PostRainSqlx(db *sql.DB, reqrain *ReqAddRain) (*Rain, string) {

	value := reqrain.Value

	var rain Rain

	lastInsertId := 0

	//sqlStatement := fmt.Sprintf("insert into 'pioggia' ('valore','data_inserimento') values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	sqlStatement := fmt.Sprintf("insert into pioggia (valore,data_inserimento) values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	log.Println(sqlStatement)

	err := db.QueryRow(sqlStatement).Scan(&lastInsertId)

	if err != nil {
		return &rain, ErrHandler(err)
	}

	sqlStatement1 := fmt.Sprintf("SELECT id,valore,data_inserimento FROM pioggia where id = %d", lastInsertId)
	rows, err := db.Query(sqlStatement1)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Rain
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		rain = p
	}
	if err != nil {
		return &rain, lang.Get("no_result")
	}
	return &rain, ""
}
