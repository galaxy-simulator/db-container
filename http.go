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
func newTreeEndpoint(width float64) {
	db := connectToDB()
	newTree(db, width)
}

// insertStarEndpoint inserts the star into the tree with the given index
func insertStarEndpoint(star structs.Star2D, index int64) {
	db := connectToDB()
	insertStar(db, star, index)
}

// deleteStarsEndpoint deletes all the rows from the stars table
func deleteStarsEndpoint() {
	db := connectToDB()
	deleteAllStars(db)
}

// deleteNodesEndpoint deletes all the rows from the nodes table
func deleteNodesEndpoint() {
	db := connectToDB()
	deleteAllNodes(db)
}