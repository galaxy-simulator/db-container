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
	"flag"
	"fmt"
	"git.darknebu.la/GalaxySimulator/structs"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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
	newTreeEndpoint(db, width)

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

	insertStarEndpoint(db, star, index)
}

func insertListHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The insertStarHandler was accessed")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	filename := r.Form.Get("filename")

	insertListEndpoint(db, filename)
}

func deleteStarsHandler(w http.ResponseWriter, r *http.Request) {
	deleteStarsEndpoint(db)
}

func deleteNodesHandler(w http.ResponseWriter, r *http.Request) {
	deleteNodesEndpoint(db)
}

func getListOfStarsGoHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The getListOfStarsGoHandler was accessed")

	listOfStars := listOfStarsGoEndpoint(db)

	for _, star := range listOfStars {
		_, _ = fmt.Fprintf(w, "%v\n", star)
	}
}

func getListOfStarsCsvHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The getListOfStarsCsvHandler was accessed")

	listOfStars := listOfStarsCsvEndpoint(db)

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

	updateTotalMassEndpoint(db, index)
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

	updateCenterOfMassEndpoint(db, index)
}

func genForestTreeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[ ] The genForestTreeHandler was accessed")

	index, _ := strconv.ParseInt(r.Form.Get("index"), 10, 64)

	tree := genForestTreeEndpoint(db, index)
	_, _ = fmt.Fprintf(w, "%s", tree)
}

func initStarsTableHandler(w http.ResponseWriter, r *http.Request) {
	initStarsTableEndpoint(db)
}

func initNodesTableHandler(w http.ResponseWriter, r *http.Request) {
	initNodesTableEndpoint(db)
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "port used to host the service")
	var dbURL string
	flag.StringVar(&dbURL, "DBURL", "postgres", "url of the database used")
	flag.Parse()

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

	router.HandleFunc("/genforesttree", genForestTreeHandler).Methods("GET")

	router.HandleFunc("/initStarsTable", initStarsTableHandler).Methods("POST")
	router.HandleFunc("/initNodesTable", initNodesTableHandler).Methods("POST")

	connStr := fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=none", dbURL)
	db, _ := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(75)

	fmt.Printf("Database Container up on localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}
