package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

// fastInsertHandler gets a tree index and a filename and tries to read insert all the stars from the file into the tree
func fastInsertJSONHandler(w http.ResponseWriter, r *http.Request) {
	// read the mux variables
	vars := mux.Vars(r)
	filename, _ := vars["filename"]

	// read the content using the given filename
	content, readErr := ioutil.ReadFile(fmt.Sprintf("/db/%s", filename))
	if readErr != nil {
		panic(readErr)
	}

	// unmarshal the file content
	tree := &structs.Node{}
	jsonUnmarshalErr := json.Unmarshal(content, tree)
	if jsonUnmarshalErr != nil {
		panic(jsonUnmarshalErr)
	}

	// append the tree to the treeArray
	treeArray = append(treeArray, tree)

	// return the treeArray index the tree was inserted into (== the length of the array)
	_, _ = fmt.Fprintln(w, len(treeArray))
}

func fastInsertListHandler(w http.ResponseWriter, r *http.Request) {
	// read the mux variables
	vars := mux.Vars(r)
	filename, _ := vars["filename"]
	treeindex, _ := strconv.ParseInt(vars["filename"], 10, 64)

	// read the content using the given filename
	content, readErr := ioutil.ReadFile(fmt.Sprintf("/home/db/%s", filename))
	if readErr != nil {
		panic(readErr)
	}

	in := string(content)
	reader := csv.NewReader(strings.NewReader(in))

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("------------------------")
			log.Println("error:")
			log.Println(record)
			log.Println(err)
			log.Println("------------------------")
		}
		fmt.Println(record)

		x, _ := strconv.ParseFloat(record[0], 64)
		y, _ := strconv.ParseFloat(record[1], 64)

		fmt.Println("___________________")
		fmt.Println(record[1])
		fmt.Println(y)
		fmt.Println("___________________")

		star := structs.NewStar2D(structs.Vec2{x, y}, structs.Vec2{0, 0}, 42)

		fmt.Printf("Star: %v", star)

		err = treeArray[treeindex].Insert(star)
		if err != nil {
			log.Println(err)
			errorCount[treeindex] += 1
		}
		fmt.Printf("Inserted %v\n", star)
		starCount[treeindex] += 1
	}
}
