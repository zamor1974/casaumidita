package models

import (
	"casaumidita/lang"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// swagger:model Humidity
type Humidity struct {
	// Id of rain value
	// in: int64
	Id int64 `json:"id"`
	// Value of Humidity
	// in: int
	Value float32 `json:"valore"`
	// Timestamp of insert
	// in: time
	DateInsert time.Time `json:"data_inserimento"`
}

type Humidities []Humidity

type ReqAddHumidity struct {
	// Value of the Humidity
	// in: int
	Value int `json:"valore" validate:"required"`
}

// swagger:parameters add Humidity
type ReqHumidityBody struct {
	// - name: body
	//  in: body
	//  description: Humidity
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddHumidity"
	//  required: true
	Body ReqAddHumidity `json:"body"`
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

func GetHumiditiesSqlx(db *sql.DB) *Humidities {
	humidities := Humidities{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM umidita order by id desc")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Humidity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		humidities = append(humidities, p)
	}
	return &humidities
}
func GetLastHumiditySqlx(db *sql.DB) *Humidities {
	humidities := Humidities{}
	rows, err := db.Query("SELECT id, valore, data_inserimento FROM umidita where id = (select max(id) from pioggia)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Humidity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		humidities = append(humidities, p)
	}
	return &humidities
}
func GetLastHourSqlx(db *sql.DB) *Humidities {
	humidities := Humidities{}

	tFine := time.Now()
	dataFine := tFine.Format("2006-01-02 15:04:05")

	tInizio := time.Now().Add(time.Duration(-1) * time.Hour)
	dataInizio := tInizio.Format("2006-01-02 15:04:05")

	sqlStatement := fmt.Sprintf("SELECT id,valore,data_inserimento FROM umidita where data_inserimento  >= '%s' AND data_inserimento <= '%s'", dataInizio, dataFine)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Humidity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			log.Fatal(err)
		}
		humidities = append(humidities, p)
	}

	if len(humidities) == 0 {
		elemento := GetLastHumiditySqlx(db)
		humidities = append(humidities, *elemento...)
	}
	return &humidities
}

// PostHumiditySqlx insert Humidity value
func PostHumiditySqlx(db *sql.DB, reqHumidity *ReqAddHumidity) (*Humidity, string) {

	value := reqHumidity.Value

	var humidity Humidity

	lastInsertId := 0

	//sqlStatement := fmt.Sprintf("insert into 'pioggia' ('valore','data_inserimento') values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	sqlStatement := fmt.Sprintf("insert into umidita (valore,data_inserimento) values (%d,CURRENT_TIMESTAMP) RETURNING id", value)
	log.Println(sqlStatement)

	err := db.QueryRow(sqlStatement).Scan(&lastInsertId)

	if err != nil {
		return &humidity, ErrHandler(err)
	}

	sqlStatement1 := fmt.Sprintf("SELECT id,valore,data_inserimento FROM umidita where id = %d", lastInsertId)
	rows, err := db.Query(sqlStatement1)

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var p Humidity
		if err := rows.Scan(&p.Id, &p.Value, &p.DateInsert); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
		humidity = p
	}
	if err != nil {
		return &humidity, lang.Get("no_result")
	}
	return &humidity, ""
}
