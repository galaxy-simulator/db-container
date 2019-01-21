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

// newTreeHandler creates a new tree
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
	width := 100.0

	// value from the user
	userWidth, _ := strconv.ParseFloat(r.Form.Get("w"), 64) // bounding box width
	log.Printf("width: %f", userWidth)

	// overwrite the default width
	if userWidth != width {
		width = userWidth
	}

	jsonData := newTree(width)

	// return the new tree as json
	_, _ = fmt.Fprintf(w, "%v", string(jsonData))
}

// newTree generates a new tree using the width it is given and returns it as json in an array of bytes
func newTree(width float64) []byte {

	// generate a new tree and add it to the treeArray
	newTree := structs.NewRoot(width)
	treeArray = append(treeArray, newTree)
	starCount = append(starCount, 0)

	// convert the tree to json format
	jsonData, jsonMarshalErr := json.Marshal(newTree)
	if jsonMarshalErr != nil {
		panic(jsonMarshalErr)
	}

	return jsonData

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

	treeInsertError := treeArray[treeindex].Insert(s1)
	if treeInsertError != nil {
		panic(fmt.Sprintf("Error inserting %v into tree %d: %v", s1, treeindex, treeInsertError))
	}

	fmt.Println("-------------")
	fmt.Println(treeArray)
	fmt.Println("-------------")

	log.Println("Done inserting the star")
	starCount[treeindex] += 1

	pushMetricsNumOfStars("http://db:80/metrics", treeindex)

	_, _ = fmt.Fprintf(w, "%d", starCount[treeindex])
}

// pushMetricsNumOfStars pushes the amount of stars in the given galaxy with the given index to the given host
// the host is (normally) the service bundling the metrics
func pushMetricsNumOfStars(host string, treeindex int64) {

	// define a post-request and send it to the given host
	requestURL := fmt.Sprintf("%s", host)
	resp, err := http.PostForm(requestURL,
		url.Values{
			"key":   {fmt.Sprintf("db_%s{nr=\"%s\"}", "stars_num", treeindex)},
			"value": {fmt.Sprintf("%d", starCount[treeindex])},
		},
	)
	if err != nil {
		fmt.Printf("Cound not make a POST request to %s", requestURL)
	}

	// close the response body
	defer resp.Body.Close()
}

// starlistHandler lists all the stars in the given tree
func starlistHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("The starlist handler was accessed")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	listofallstars := treeArray[treeindex].GetAllStars()
	log.Printf("listofallstars: %v", listofallstars)
	// listofallstars: [{{-42 10} {0 0} 100} {{10 10} {0 0} 100}]

	// convert the list of all stars to json
	jsonlistofallstars, jsonMarshalErr := json.Marshal(listofallstars)
	if jsonMarshalErr != nil {
		panic(jsonMarshalErr)
	}

	log.Printf("jsonlistofallstars: %v", string(jsonlistofallstars))

	_, _ = fmt.Fprintln(w, string(jsonlistofallstars))
	log.Println("Done")
}

// dumptreeHandler dumps the requested tree
func dumptreeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("The dumptree endpoint was accessed.\n")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	// Convert the data to json
	jsonData, jsonMarshalerError := json.Marshal(treeArray[treeindex])
	if jsonMarshalerError != nil {
		panic(jsonMarshalerError)
	}

	// print the jsonData to the ResponseWriter
	_, printTreeErr := fmt.Fprintf(w, "%v\n", string(jsonData))
	if printTreeErr != nil {
		panic(printTreeErr)
	}
}

// updateCenterOfMassHandler updates the center of mass in each node in tree with the given index
func updateCenterOfMassHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating the center of mass")
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	treeArray[treeindex].CalcCenterOfMass()
}

// updateTotalMassHandler updates the total mass in each node in the tree with the given index
func updateTotalMassHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Updating the total mass")
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	treeArray[treeindex].CalcTotalMass()
}

// metricHandler prints all the metrics to the ResponseWriter
func metricHandler(w http.ResponseWriter, r *http.Request) {
	var metricsString string
	metricsString += fmt.Sprintf("nr_galaxies %d\n", len(treeArray))

	for i := 0; i < len(starCount); i++ {
		metricsString += fmt.Sprintf("galaxy_star_count{galaxy_nr=\"%d\"} %d\n", i, starCount[i])
	}

	log.Println(metricsString)
	_, _ = fmt.Fprintf(w, metricsString)
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

	fmt.Println("Database Container up")
	log.Fatal(http.ListenAndServe(":80", router))
}
