package main

import (
	"app/database"
	"encoding/json"
	"strconv"
	//"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type M2TParams struct {
	MId int64 `json:"mId"`
	TId int64 `json:"tagId"`
}

var db *database.Database

func routes(_db *database.Database) {
	db = _db

	r := mux.NewRouter()
	r.HandleFunc("/stats/day/{date:[0-9]{4}-[0-9]{2}-[0-9]{1,2}}", totalsForDayHandler)
	r.HandleFunc("/stats/unmatched", totalsUnmatchedHandler)
	r.HandleFunc("/stats/tags", totalsByTagHandler)
	r.HandleFunc("/tags", TagIndex).Methods("GET")
	r.HandleFunc("/tags/{name}", TagGet).Methods("GET", "POST")
	r.HandleFunc("/tags", TagCreate).Methods("POST")

	r.HandleFunc("/matchers", MatcherIndex).Methods("GET")
	r.HandleFunc("/matchers/{id}", MatcherGet).Methods("GET", "POST")
	r.HandleFunc("/matchers/{id}", MatcherDelete).Methods("DELETE")
	r.HandleFunc("/matchers", MatcherCreate).Methods("POST")

	r.HandleFunc("/me2tags", Matcher2TagCreate).Methods("POST")
	r.HandleFunc("/me2tags", Matcher2TagDelete).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	//r.PathPrefix("/").Handler(http.FileServer(
	//&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "static"}))
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

	wr.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(wr).Encode(db.TotalsForDay(date)); err != nil {
		panic(err)
	}
}

func totalsUnmatchedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	events_all := database.GetDB().EventsAllChannel()
	events_filtered := EventRecordFilterUnmatched(events_all)

	out := make([]database.EventRecord, 0)
	for record := range events_filtered {
		out = append(out, record)
	}

	j, _ := json.Marshal(out)
	w.Write(j)
}

func TagIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(database.GetDB().GetTags())
	w.Write(j)
}

func TagCreate(w http.ResponseWriter, r *http.Request) {
	tag := new(database.Tag)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tag)
	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(database.GetDB().CreateTag(tag))
	w.Write(j)
}

func TagGet(w http.ResponseWriter, r *http.Request) {
}

func MatcherCreate(w http.ResponseWriter, r *http.Request) {
	matcher := new(database.MatchExpression)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&matcher)

	AppendExpression(matcher)

	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(database.GetDB().CreateMatchExpression(matcher))
	w.Write(j)
}

func MatcherIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(database.GetDB().GetMatchExpressions())
	w.Write(j)
}

func MatcherGet(w http.ResponseWriter, r *http.Request) {

}

func MatcherDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	database.GetDB().DeleteMatchExpressionById(id)

	ReloadExpressions()
}

func parseM2TParams(r *http.Request) *M2TParams {
	decoder := json.NewDecoder(r.Body)
	var params M2TParams
	err := decoder.Decode(&params)
	if err != nil {
		panic(err)
	}
	return &params
}

func Matcher2TagCreate(w http.ResponseWriter, r *http.Request) {
	params := parseM2TParams(r)
	database.GetDB().M2TCreate(params.MId, params.TId)
}

func Matcher2TagDelete(w http.ResponseWriter, r *http.Request) {
	params := parseM2TParams(r)
	database.GetDB().M2TDestroy(params.MId, params.TId)
}

func totalsByTagHandler(w http.ResponseWriter, r *http.Request) {
	events_all := database.GetDB().EventsAllChannel()
	totals := GetTotalsByTag(events_all)

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(totals)
	w.Write(j)
}
