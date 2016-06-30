package main

import (
	"app/database"
	"encoding/json"
	//"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var db *database.Database

func routes(_db *database.Database) {
	db = _db

	r := mux.NewRouter()
	r.HandleFunc("/stats/day/{date:[0-9]{4}-[0-9]{2}-[0-9]{1,2}}", totalsForDayHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8228", nil))
}

func totalsForDayHandler(wr http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	datestring := vars["date"]

	date, err := time.Parse("2006-01-02", datestring)
	if err != nil {
		panic(err)
	}

	records := db.TotalsForDay(date)

	s := make([]database.TotalStatisticsRecord, 0)
	for record := range records {
		s = append(s, record)
	}

	wr.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(wr).Encode(s); err != nil {
		panic(err)
	}

}
