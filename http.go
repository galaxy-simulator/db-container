// http.go bundles the http endpoint definitions
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
	"git.darknebu.la/GalaxySimulator/db_actions"
	"git.darknebu.la/GalaxySimulator/structs"
)

// IndexEndpoint gives a basic overview over the api
func indexEndpoint() string {
	indexString := `Galaxy Simulator Database

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

	return indexString
}

// newTree creates a new tree
func newTreeEndpoint(db *sql.DB, width float64) {
	db_actions.NewTree(db, width)
}

// insertStarEndpoint inserts the star into the tree with the given index
func insertStarEndpoint(db *sql.DB, star structs.Star2D, index int64) {
	db_actions.InsertStar(db, star, index)
}

// insertListEndpoint inserts the star into the tree with the given index
func insertListEndpoint(db *sql.DB, filename string) {
	db_actions.InsertList(db, filename)
}

// deleteStarsEndpoint deletes all the rows from the stars table
func deleteStarsEndpoint(db *sql.DB) {
	db_actions.DeleteAllStars(db)
}

// deleteNodesEndpoint deletes all the rows from the nodes table
func deleteNodesEndpoint(db *sql.DB) {
	db_actions.DeleteAllNodes(db)
}

func listOfStarsGoEndpoint(db *sql.DB) []structs.Star2D {
	listOfStars := db_actions.GetListOfStarsGo(db)
	return listOfStars
}

func listOfStarsCsvEndpoint(db *sql.DB) []string {
	listOfStars := db_actions.GetListOfStarsCsv(db)
	return listOfStars
}

func updateTotalMassEndpoint(db *sql.DB, index int64) {
	db_actions.UpdateTotalMass(db, index)
}

func updateCenterOfMassEndpoint(db *sql.DB, index int64) {
	db_actions.UpdateCenterOfMass(db, index)
}

func genForestTreeEndpoint(db *sql.DB, index int64) string {
	return db_actions.GenForestTree(db, index)
}
