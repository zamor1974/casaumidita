package main

import (
	"casaumidita/config"
	"casaumidita/controllers"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
	})
	dbsqlx := config.ConnectDBSqlx()
	hsqlx := controllers.NewBaseHandlerSqlx(dbsqlx)

	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	// documentation for developers
	opts := middleware.SwaggerUIOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// documentation for share
	opts1 := middleware.RedocOpts{SpecURL: "/swagger.yaml", Path: "docs"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/docs", sh1)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{\"hello\": \"world\"}"))
	})

	unidita := r.PathPrefix("/").Subrouter()
	unidita.HandleFunc("/humidity", hsqlx.PostHumiditySqlx).Methods("POST")
	unidita.HandleFunc("/humidities", hsqlx.GetHumiditiesSqlx).Methods("GET")
	unidita.HandleFunc("/lasthour", hsqlx.GetLastHourSqlx).Methods("GET")

	http.Handle("/", r)
	s := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", "", "5556"),
		Handler: cors.Default().Handler(r),
	}
	s.ListenAndServe()
}
