package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	treeArray []*structs.Node
	starCount []int
)

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	infostring := `Galaxy Simulator Database

API:
	/ GET
	/new POST w float64
	/insert/{treeindex} POST x float64, y float64, vx float64, vy float64, m float64
	/starlist/{treeindex}
	/printall GET
	/metrics GET
`
	_, _ = fmt.Fprintf(w, infostring)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/new", newTreeHandler).Methods("POST")
	router.HandleFunc("/printall", printAllHandler).Methods("GET")
	router.HandleFunc("/insert/{treeindex}", insertStarHandler).Methods("POST")
	router.HandleFunc("/starlist/{treeindex}", starlistHandler).Methods("GET")
	router.HandleFunc("/dumptree/{treeindex}", dumptreeHandler).Methods("GET")
	router.HandleFunc("/updatetotalmass/{treeindex}", updateTotalMassHandler).Methods("GET")
	router.HandleFunc("/updatecenterofmass/{treeindex}", updateCenterOfMassHandler).Methods("GET")
	router.HandleFunc("/metrics", metricHandler).Methods("GET")
	router.HandleFunc("/export/{treeindex}", exportHandler).Methods("POST")
	router.HandleFunc("/fastinsert/{filename}", fastInsertHandler).Methods("POST")

	fmt.Println("Database Container up")
	log.Fatal(http.ListenAndServe(":80", router))
}
