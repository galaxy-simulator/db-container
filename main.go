package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	treeArray  []*structs.Node
	starCount  []int
	errorCount []int
)

func index() string {
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

	- /fastinsertjson/{filename} ("GET")
	- /fastinsertlist/{filename} ("GET")

	- /readdir ("GET")
`
	return infostring
}

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, index())
}

// readdirHandler reads the content of a given directory and prints it
func readdirHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dirname := vars["dirname"]
	log.Printf("Reading from %s", dirname)

	files, err := ioutil.ReadDir(fmt.Sprintf("./%s", dirname))
	log.Println(files)
	log.Println(err)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range files {
		fmt.Println(f.Name())
		_, _ = fmt.Fprintf(w, "%v", f.Name())
	}
}

// return the amount of galaxies currently present
func nrofgalaxiesHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "%d", len(treeArray))
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
	router.HandleFunc("/nrofgalaxies", nrofgalaxiesHandler).Methods("GET")

	router.HandleFunc("/fastinsertjson/{filename}", fastInsertJSONHandler).Methods("GET")
	router.HandleFunc("/fastinsertlist/{filename}/{treeindex}", fastInsertListHandler).Methods("GET")

	router.HandleFunc("/readdir/{dirname}", readdirHandler).Methods("GET")

	fmt.Println("Database Container up")
	log.Fatal(http.ListenAndServe(":80", router))
}
