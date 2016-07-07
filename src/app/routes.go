package main

import (
	"encoding/json"
	"fmt"
	//"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func routes(_db *Database) {
	db = _db

	r := mux.NewRouter()
	r.HandleFunc("/stats/day/{date:[0-9]{4}-[0-9]{2}-[0-9]{1,2}}", totalsForDayHandler)
	r.HandleFunc("/stats/unmatched", totalsUnmatchedHandler)
	r.HandleFunc("/stats/tags", totalsByTagHandler)
	r.HandleFunc("/stats/tags/tree", totalsByTagTreeHandler)
	r.HandleFunc("/stats/tags/{parentId}", totalsByTagHandler)

	r.HandleFunc("/tags", TagIndex).Methods("GET")
	r.HandleFunc("/tags/{name}", TagGet).Methods("GET", "POST")
	r.HandleFunc("/tags", TagCreate).Methods("POST")
	r.HandleFunc("/tags/{id}", tagDeleteHandler).Methods("DELETE")

	r.HandleFunc("/matchers", MatcherIndex).Methods("GET")
	r.HandleFunc("/matchers/{id}", MatcherGet).Methods("GET", "POST")
	r.HandleFunc("/matchers/{id}", MatcherDelete).Methods("DELETE")
	r.HandleFunc("/matchers", MatcherCreate).Methods("POST")

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

	events_all := GetDB().EventsAllChannel()
	events_filtered := EventRecordFilterUnmatched(events_all)

	out := make([]EventRecord, 0)
	for record := range events_filtered {
		out = append(out, record)
	}

	j, _ := json.Marshal(out)
	w.Write(j)
}

func TagIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	j, _ := json.Marshal(GetDB().GetTags())
	w.Write(j)
}

func TagCreate(w http.ResponseWriter, r *http.Request) {
	tag := new(Tag)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&tag)
	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(GetDB().CreateTag(tag))
	w.Write(j)
}

func TagGet(w http.ResponseWriter, r *http.Request) {
}

func tagDeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	GetDB().DeleteTagById(id)
}

func MatcherCreate(w http.ResponseWriter, r *http.Request) {
	matcher := new(MatchExpression)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&matcher)

	AppendExpression(matcher)

	if err != nil {
		panic(err)
	}

	j, _ := json.Marshal(GetDB().CreateMatchExpression(matcher))
	w.Write(j)
}

func MatcherIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(GetDB().GetMatchExpressions())
	w.Write(j)
}

func MatcherGet(w http.ResponseWriter, r *http.Request) {

}

func MatcherDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.ParseInt(vars["id"], 10, 64)
	GetDB().DeleteMatchExpressionById(id)

	ReloadExpressions()
}

func totalsByTagHandler(w http.ResponseWriter, r *http.Request) {
	events_all := GetDB().EventsAllChannel()

	vars := mux.Vars(r)
	parentId, _ := strconv.ParseInt(vars["parentId"], 10, 64)

	var matchers *[]CompiledMatchExpression

	if parentId > -1 {
		matchers = GetMatchersWithParentId(int(parentId))
	} else {
		matchers = GetMatchers()
	}

	totals := GetTotalsByTag(events_all, matchers, true)

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(totals)
	w.Write(j)
}

func totalsByTagTreeHandler(w http.ResponseWriter, r *http.Request) {
	events_all := GetDB().EventsAllChannel()
	var matchers *[]CompiledMatchExpression

	matchers = GetMatchers()
	totals := GetTotalsByTag(events_all, matchers, true)
	//tree := BuildTagTotalsTree(totals)
	//tree := TreeifyTagTotals(totals)
	//for _, item := range tree {
	//fmt.Printf("%s\n", item)
	//}
	tree := BuildTagTotalsTree(totals)
	fmt.Printf("tree %s\n", tree)

	w.Header().Set("Content-Type", "application/json")
	j, _ := json.Marshal(tree)
	w.Write(j)
}
