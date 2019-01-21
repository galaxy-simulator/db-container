package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

// fastInsertHandler gets a tree index and a filename and tries to read
func fastInsertHandler(w http.ResponseWriter, r *http.Request) {
	// read the mux variables
	vars := mux.Vars(r)
	filename, _ := vars["filename"]

	// read the content using the given filename
	content, readErr := ioutil.ReadFile(filename)
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
