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
	- / ("GET")
	- /new ("POST")
	- /printall ("GET")
	- /insert/{treeindex} ("POST")
	- /starlist/{treeindex} ("GET")
	- /dumptree/{treeindex} ("GET")
	- /updatetotalmass/{treeindex} ("GET")
	- /updatecenterofmass/{treeindex} ("GET")
	- /metrics ("GET")
	- /export/{treeindex} ("POST")
	- /fastinsert/{filename} ("POST")
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
