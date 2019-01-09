package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	treeArray []structs.Quadtree
)

// Get a subtree by searching a given element and returning its children recursively
func getSubtreeHandler(w http.ResponseWriter, r *http.Request) {
	// set the content type to json (looks fancy in firefox :D)
	w.Header().Set("Content-Type", "application/json")

	// parse the mux variables
	vars := mux.Vars(r)
	getIndex, strconvErr := strconv.ParseInt(vars["treeindex"], 10, 0)
	if strconvErr != nil {
		panic(strconvErr)
	}
	log.Println(getIndex)

	// Convert the data to json
	jsonData, jsonMarshalerError := json.Marshal(treeArray[getIndex])
	if jsonMarshalerError != nil {
		panic(jsonMarshalerError)
	}

	// print the jsonData to the ResponseWriter
	_, printTreeErr := fmt.Fprintf(w, "%v\n", string(jsonData))
	if printTreeErr != nil {
		panic(printTreeErr)
	}
	log.Printf("The getSubtree endpoint was accessed.\n")
}

// newTreeHandler creates a new tree and adds ot the the treeArray
func newTreeHandler(w http.ResponseWriter, r *http.Request) {
	// set the content type to json (looks fancy in firefox :D)
	w.Header().Set("Content-Type", "application/json")

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// default values
	x := 0.0
	y := 0.0
	width := 0.0

	// values from the user
	xTmp, _ := strconv.ParseFloat(r.Form.Get("x"), 64)     // x
	yTmp, _ := strconv.ParseFloat(r.Form.Get("y"), 64)     // y
	widthTmp, _ := strconv.ParseFloat(r.Form.Get("w"), 64) // bounding box width

	// assign the values
	if xTmp != 0 {
		x = xTmp
	}
	if yTmp != 0 {
		y = yTmp
	}
	if widthTmp != 0 {
		width = widthTmp
	}

	// generate a new tree and add it to the treeArray
	newTree := structs.NewQuadtree(structs.BoundingBox{
		Center: structs.Vec2{
			X: x,
			Y: y,
		},
		Width: width,
	})

	log.Println(newTree.Boundary)

	treeArray = append(treeArray, *newTree)

	// convert the tree to json format
	jsonData, jsonMarshalErr := json.Marshal(newTree)
	if jsonMarshalErr != nil {
		panic(jsonMarshalErr)
	}

	// return the new tree as json
	_, _ = fmt.Fprintf(w, "%v", string(jsonData))

	log.Printf("The newTree endpoint was accessed.\n")
}

// printAllHandler prints all the trees in the treeArray	router.HandleFunc("/printall", printAllHandler).Methods("GET")
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

func generatePrintTree(quadtree structs.Quadtree) string {
	returnString := "["
	fmt.Printf("[")
	for i := 0; i < 4; i++ {
		if quadtree.Quadrants[i] != nil {
			returnString += generatePrintTree(*quadtree.Quadrants[i])
		}
	}
	returnString += "]"
	fmt.Printf("]")
	return returnString
}

func printTreeHandler(w http.ResponseWriter, r *http.Request) {
	returnString := generatePrintTree(treeArray[0])
	_, _ = fmt.Fprintln(w, returnString)
}

// this insert handler inserts a given star using http queries
func insertHandler(w http.ResponseWriter, r *http.Request) {
	// get the tree id in which the star should be inserted
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)
	_, _ = fmt.Fprintln(w, treeindex)

	// get the star by parsing http-post parameters
	errParseForm := r.ParseForm() // parse the POST form
	if errParseForm != nil {      // handle errors
		panic(errParseForm)
	}

	// parse the values from the post parameters	router.HandleFunc("/printall", printAllHandler).Methods("GET")
	x, _ := strconv.ParseFloat(r.Form.Get("x"), 64)
	y, _ := strconv.ParseFloat(r.Form.Get("y"), 64)
	vx, _ := strconv.ParseFloat(r.Form.Get("vx"), 64)
	vy, _ := strconv.ParseFloat(r.Form.Get("vy"), 64)
	m, _ := strconv.ParseFloat(r.Form.Get("m"), 64)

	log.Printf("[---] Inserting star into the tree")

	// build the star that should be inserted
	newStar := structs.Star2D{
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

	treeArray[treeindex].NewInsert(newStar)
}

// Simple index Handler
// TODO: Display some kind of help
// TODO: Insert an api-documentation
func indexHandler(w http.ResponseWriter, r *http.Request) {
	var _, _ = fmt.Fprintf(w, "Hello from the db-container!")
	log.Printf("The indexHandler was accessed.")
	var _, _ = fmt.Fprintln(w, "Insert a star using > $ curl --data \"x=250&y=250&vx=0.1&vy=0.2&m=3\" http://localhost:8123/insert/0")
}

// drawGalaxyHandler draws the galaxy and returns an image of it
func drawGalaxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("The drawTreeHandler was accessed.")

	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)
	log.Println(treeindex)

	if treeArray[treeindex] != (structs.Quadtree{}) {
		log.Println(treeArray[treeindex])
		treeArray[treeindex].DrawGalaxy("/public/quadtree.png")
	}

	http.ServeFile(w, r, "/public/quadtree.png")
}

// drawTreeHandler draws a tree of the galaxy and returns an image of it
func drawTreeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("The drawTreeHandler was accessed.")

	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)
	log.Println(treeindex)

	if treeArray[treeindex] != (structs.Quadtree{}) {
		log.Println(treeArray[treeindex])
		latex := treeArray[treeindex].DrawTree()
		_, _ = fmt.Fprintf(w, "%s", latex)
	} else {
		_, _ = fmt.Fprintln(w, "error")
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/get/{treeindex}", getSubtreeHandler).Methods("GET")
	router.HandleFunc("/new", newTreeHandler).Methods("POST")
	router.HandleFunc("/insert/{treeindex}", insertHandler).Methods("POST")
	router.HandleFunc("/printall", printAllHandler).Methods("GET")
	router.HandleFunc("/printtree", printTreeHandler).Methods("GET")
	router.HandleFunc("/drawgalaxy/{treeindex}", drawGalaxyHandler).Methods("GET")
	router.HandleFunc("/drawtree/{treeindex}", drawTreeHandler).Methods("GET")

	log.Println("Serving the database on port 8092 (This is for local testing only, remove when done!)")
	log.Fatal(http.ListenAndServe(":8042", router))
}
