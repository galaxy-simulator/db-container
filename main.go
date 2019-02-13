// main.go purpose is to build the interaction layer in between the http endpoints and the http server
// Copyright (C) 2019 Emile Hansmaennel
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"database/sql"
	"fmt"
	"git.darknebu.la/GalaxySimulator/db-container/db_actions"
	"git.darknebu.la/GalaxySimulator/structs"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	//"git.darknebu.la/GalaxySimulator/structs"
)

var (
	db *sql.DB
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The indexHandler was accessed")
	_, _ = fmt.Fprintf(w, indexEndpoint())
}

func newTreeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The newTreeHandler was accessed")

	// get the width of the most outer box of the tree by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// parse the width
	width, _ := strconv.ParseFloat(r.Form.Get("w"), 64)

	// create a new tree in the database width the given width
	newTreeEndpoint(width)

	_, _ = fmt.Fprintf(w, "OK\n")
}

func insertStarHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The insertStarHandler was accessed")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// get the star coordinates
	x, _ := strconv.ParseFloat(r.Form.Get("x"), 64)
	y, _ := strconv.ParseFloat(r.Form.Get("y"), 64)
	vx, _ := strconv.ParseFloat(r.Form.Get("vx"), 64)
	vy, _ := strconv.ParseFloat(r.Form.Get("vy"), 64)
	m, _ := strconv.ParseFloat(r.Form.Get("m"), 64)

	// get the tree into which the star should be inserted into
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)

	// build a star
	star := structs.Star2D{
		C: structs.Vec2{
			X: x,
			Y: y,
		},
		V: structs.Vec2{
			X: vx,
			Y: vy,
		},
		M: m,
	}

	insertStarEndpoint(star, index)
}

func insertListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The insertStarHandler was accessed")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	filename := r.Form.Get("filename")

	insertListEndpoint(filename)
}

func deleteStarsHandler(w http.ResponseWriter, r *http.Request) {
	deleteStarsEndpoint()
}

func deleteNodesHandler(w http.ResponseWriter, r *http.Request) {
	deleteNodesEndpoint()
}

func getListOfStarsGoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The getListOfStarsGoHandler was accessed")

	listOfStars := listOfStarsGoEndpoint()

	for _, star := range listOfStars {
		_, _ = fmt.Fprintf(w, "%v\n", star)
	}
}

func getListOfStarsCsvHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The getListOfStarsCsvHandler was accessed")

	listOfStars := listOfStarsCsvEndpoint()

	for _, star := range listOfStars {
		_, _ = fmt.Fprintf(w, "%v\n", star)
	}
}

func updateTotalMassHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The updateTotalMassHandler was accessed")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// get the tree into which the star should be inserted into
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)

	updateTotalMassEndpoint(index)
}

func updateCenterOfMassHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The updateCenterOfMassHandler was accessed")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// get the tree into which the star should be inserted into
	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)

	updateCenterOfMassEndpoint(index)
}

func genForestTreeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The genForestTreeHandler was accessed")

	vars := mux.Vars(r)
	treeindex, parseIntErr := strconv.ParseInt(vars["treeindex"], 10, 64)
	if parseIntErr != nil {
		panic(parseIntErr)
	}

	tree := genForestTreeEndpoint(treeindex)
	_, _ = fmt.Fprintf(w, "%s", tree)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/new", newTreeHandler).Methods("POST")
	router.HandleFunc("/deleteStars", deleteStarsHandler).Methods("POST")
	router.HandleFunc("/deleteNodes", deleteNodesHandler).Methods("POST")
	router.HandleFunc("/starlist/go", getListOfStarsGoHandler).Methods("GET")
	router.HandleFunc("/starlist/csv", getListOfStarsCsvHandler).Methods("GET")

	router.HandleFunc("/insertStar", insertStarHandler).Methods("POST")
	router.HandleFunc("/insertList", insertListHandler).Methods("POST")

	router.HandleFunc("/updatetotalmass", updateTotalMassHandler).Methods("POST")
	router.HandleFunc("/updatecenterofmass", updateCenterOfMassHandler).Methods("POST")

	router.HandleFunc("/genforesttree/{treeindex}", genForestTreeHandler).Methods("GET")

	//router.HandleFunc("/metrics", metricHandler).Methods("GET")
	//router.HandleFunc("/export/{treeindex}", exportHandler).Methods("POST")
	//router.HandleFunc("/nrofgalaxies", nrofgalaxiesHandler).Methods("GET")

	db = db_actions.ConnectToDB()
	db.SetMaxOpenConns(75)

	fmt.Println("Database Container up on port 8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}
