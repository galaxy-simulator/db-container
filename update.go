package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
