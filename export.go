package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// export exports all the trees
func export(treeindex int64) error {
	// Convert the data to json
	jsonData, jsonMarshalerError := json.Marshal(treeArray[treeindex])
	if jsonMarshalerError != nil {
		panic(jsonMarshalerError)
	}

	// write the json formatted byte data to a file
	err := ioutil.WriteFile(fmt.Sprintf("/exports/tree_%d.json", treeindex), jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}

// export the selected tree to the selected file
func exportHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	err := export(treeindex)
	if err != nil {
		panic(err)
	}

	_, _ = fmt.Fprintf(w, "Exportet Tree %d", treeindex)
}
