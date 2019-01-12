package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	treeArray []*structs.Node
)

// indexHandler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	infostring := `Galaxy Simulator Database

API: /api/v1/...
	.../new
	.../insert{treeindex}`
	_, _ = fmt.Fprintf(w, infostring)
}

// newTreeHandler creates a new tree and adds ot the the treeArray
func newTreeHandler(w http.ResponseWriter, r *http.Request) {
	// set the content type to json (looks fancy in firefox :D)
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("Creating a new tree")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// default values
	width := 0.0

	// value from the user
	widthTmp, _ := strconv.ParseFloat(r.Form.Get("w"), 64) // bounding box width
	log.Printf("width: %f", widthTmp)

	if widthTmp != 0 {
		width = widthTmp
	}

	// generate a new tree and add it to the treeArray
	newTree := structs.NewRoot(width)
	treeArray = append(treeArray, newTree)

	// convert the tree to json format
	jsonData, jsonMarshalErr := json.Marshal(newTree)
	if jsonMarshalErr != nil {
		panic(jsonMarshalErr)
	}

	// return the new tree as json
	_, _ = fmt.Fprintf(w, "%v", string(jsonData))

	log.Printf("The newTree endpoint was accessed.\n")
}

// printAllHandler prints all the trees in the treeArray
func printAllHandler(w http.ResponseWriter, r *http.Request) { // set the content type to json (looks fancy in firefox :D)
	w.Header().Set("Content-Type", "application/json")

	// Convert the data to json
	jsonData, jsonMarshalerError := json.Marshal(treeArray)
	if jsonMarshalerError != nil {
		panic(jsonMarshalerError)
	}

	// print the jsonData to the ResponseWriter
	_, printTreeErr := fmt.Fprintf(w, "%v\n", string(jsonData))
	if printTreeErr != nil {
		panic(printTreeErr)
	}

	log.Println(treeArray[0])
	log.Printf("The printAll endpoint was accessed.\n")
}

// insertStarHandler inserts a star into the given tree
func insertStarHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("The insert handler was accessed")

	// get the treeindex in which the star should be inserted into
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

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

	log.Printf("treeindex: %d", treeindex)
	log.Printf("x: %f", x)
	log.Printf("y: %f", y)
	log.Printf("vx: %f", vx)
	log.Printf("vy: %f", vy)
	log.Printf("m: %f", m)

	s1 := structs.Star2D{
		C: structs.Vec2{x, y},
		V: structs.Vec2{vx, vy},
		M: m,
	}

	log.Printf("s1: %v", s1)

	treeArray[treeindex].Insert(s1)

	fmt.Println("-------------")
	fmt.Println(treeArray)
	fmt.Println("-------------")

	log.Println("Done inserting the star")
}

// calculate the forces acting inbetween all the stars
func calcallHandler(w http.ResponseWriter, r *http.Request) {
	// iterate over all the stars and make a POST request to the simulator with the star

	// get the treeindex
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	listOfStars := treeArray[treeindex].GetAllStars()

	for _, star := range listOfStars {
		// http post request to the simulator traefik with the star in the form

		fmt.Println(star)

		apiurl := "simu.docker.localhost"

		response, err := http.PostForm(apiurl, url.Values{
			"x":  {fmt.Sprintf("%f", star.C.X)},
			"y":  {fmt.Sprintf("%f", star.C.Y)},
			"vx": {fmt.Sprintf("%f", star.V.X)},
			"vy": {fmt.Sprintf("%f", star.V.X)},
			"m":  {fmt.Sprintf("%f", star.M)},
		})
		if err != nil {
			panic(err)
		}

		fmt.Println(response)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/api/v1/new", newTreeHandler).Methods("POST")
	router.HandleFunc("/api/v1/printall", printAllHandler).Methods("GET")
	router.HandleFunc("/api/v1/insert/{treeindex}", insertStarHandler).Methods("POST")
	router.HandleFunc("/api/v1/calcall/{treeindex}", calcallHandler).Methods("GET")

	log.Println("Serving the database on port 8043: This is for local testing only, remove when done!)")
	log.Fatal(http.ListenAndServe(":8043", router))
}
