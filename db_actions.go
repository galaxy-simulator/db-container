// db_actions defines actions on the database
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
	"fmt"
	"git.darknebu.la/GalaxySimulator/structs"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

const (
	DB_USER    = "postgres"
	DB_NAME    = "postgres"
	DB_SSLMODE = "disable"
)

// connectToDB returns a pointer to an sql database writing to the database
func connectToDB() *sql.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=%s", DB_USER, DB_NAME, DB_SSLMODE)
	db := dbConnect(connStr)
	return db
}

// dbConnect connects to a PostgreSQL database
func dbConnect(connStr string) *sql.DB {
	// connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("[E] connection: %v", err)
	}

	return db
}

func newTree(db *sql.DB, width float64) {

	// get the current max root id
	query := fmt.Sprintf("SELECT COALESCE(max(root_id), 0) FROM nodes")
	var currentMaxRootID int64
	err := db.QueryRow(query).Scan(&currentMaxRootID)
	if err != nil {
		log.Fatalf("[E] max root id query: %v", err)
	}

	// build the query creating a new node
	query = fmt.Sprintf("INSERT INTO nodes (box_width, root_id, box_center, depth, isleaf) VALUES (%f, %d, '{0, 0}', 0, TRUE)", width, currentMaxRootID+1)

	// execute the query
	_, err = db.Query(query)
	if err != nil {
		log.Fatalf("[E] insert new node query: %v", err)
	}
}

func insertStar(db *sql.DB, star structs.Star2D, index int64) {
	// insert the star into the stars table
	log.Println("[ ] Inserting the star into the stars table")
	starID := insertIntoStars(star, db)
	log.Println("[ ] Done")

	// get the root node id
	query := fmt.Sprintf("SELECT node_id FROM nodes WHERE root_id=%d", index)
	var id int64
	err := db.QueryRow(query).Scan(&id)
	if err != nil {
		log.Fatalf("[E] Get root node id query: %v", err)
	}

	// insert the star into the tree (using it's ID) starting at the root
	log.Println("[ ] Inserting the star into the tree")
	insertIntoTree(starID, db, id)
	log.Println("[ ] Done")
}

// insertIntoStars inserts the given star into the stars table
func insertIntoStars(star structs.Star2D, db *sql.DB) int64 {
	// unpack the star
	x := star.C.X
	y := star.C.Y
	vx := star.V.X
	vy := star.V.Y
	m := star.M

	// build the request query
	query := fmt.Sprintf("INSERT INTO stars (x, y, vx, vy, m) VALUES (%f, %f, %f, %f, %f) RETURNING star_id", x, y, vx, vy, m)

	// execute the query
	var starID int64
	err := db.QueryRow(query).Scan(&starID)
	if err != nil {
		log.Fatalf("[E] insert query: %v", err)
	}

	return starID
}

func insertIntoTree(starID int64, db *sql.DB, nodeID int64) {
	// There exist four cases:
	//                    | Contains a Star | Does not Contain a Star |
	// ------------------ + --------------- + ----------------------- + ---
	// Node is a Leaf     | Impossible      | insert into node        |
	//                    |                 | subdivide               |
	// ------------------ + --------------- + ----------------------- + ---
	// Node is not a Leaf | insert preexist | insert into the subtree |
	//                    | insert new      |                         |
	// ------------------ + --------------- + ----------------------- + ---
	//                    |                 |                         |

	// get the node with the given nodeID
	// find out if the node contains a star or not
	containsStar := containsStar(db, nodeID)

	// find out if the node is a leaf
	isLeaf := isLeaf(db, nodeID)

	// if the node is a leaf and contains a star
	// throw error
	if isLeaf == true && containsStar == true {
		log.Fatalf("[E] Node is a leaf and contains a star -> impossible")
	}

	// if the node is a leaf and does not contain a star
	// insert the star into the node and subdivide it
	if isLeaf == true && containsStar == false {
		log.Println("[~] isLeaf == true && containsStar == false")
		directInsert(starID, db, nodeID)
		subdivide(db, nodeID)
	}

	// if the node is not a leaf and contains a star
	// insert the preexisting star into the correct subtree
	// insert the new star into the subtree
	if isLeaf == false && containsStar == true {
		log.Println("[~] isLeaf == false && containsStar == true")

		// Stage 1: Inserting the blocking star
		log.Println("[i] Getting the blocking-Star-ID")
		blockingStar := getBlockingStar(db, nodeID)
		log.Println("[i] Done")

		log.Println("[i] Getting the quadrant the blocking star is inside of")
		blockingStarQuadrant := quadrant(db, blockingStar, nodeID)
		log.Println("[i] Done")

		log.Println("[i] Getting the quadrant node id")
		quadrantNodeID := getQuadrantNodeID(db, nodeID, blockingStarQuadrant)
		log.Println("[i] Done")

		log.Printf("[i] Inserting the star into the new quadrant %d (%d)", quadrantNodeID, blockingStarQuadrant)
		insertIntoTree(starID, db, quadrantNodeID)
		log.Println("[i] Done")

		removeStarFromNode(db, nodeID)

		log.Println("[i] Getting the blocking-Star-ID")
		star := getStar(db, starID)
		log.Println("[i] Done")

		// Stage 2: Inserting the star that should originally be inserted
		log.Println("[i] Getting the quadrant the star is inside of")
		starQuadrant := quadrant(db, star, nodeID)
		log.Println("[i] Done")

		log.Println("[i] Getting the quadrant node id")
		quadrantNodeID = getQuadrantNodeID(db, nodeID, starQuadrant)
		log.Println("[i] Done")

		log.Printf("[i] Inserting the star into the new quadrant %d (%d)", quadrantNodeID, starQuadrant)
		insertIntoTree(starID, db, quadrantNodeID)
		log.Println("[i] Done")
	}

	// if the node is not a leaf and does not contain a star
	// insert the new star into the subtree
	if isLeaf == false && containsStar == false {
		log.Println("[~] isLeaf == false && containsStar == false")
		directInsert(starID, db, nodeID)
	}
}

// containsStar returns true if the node with the given id contains a star
// and returns false if not.
func containsStar(db *sql.DB, id int64) bool {
	var starID int64

	query := fmt.Sprintf("SELECT star_id FROM nodes WHERE node_id=%d", id)
	err := db.QueryRow(query).Scan(&starID)
	if err != nil {
		log.Fatalf("[E] containsStar query: %v", err)
	}

	if starID != 0 {
		return true
	}

	return false
}

// isLeaf returns true if the node with the given id is a leaf
func isLeaf(db *sql.DB, nodeID int64) bool {
	var isLeaf bool

	query := fmt.Sprintf("SELECT COALESCE(isleaf, FALSE) FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&isLeaf)
	if err != nil {
		log.Fatalf("[E] isLeaf query: %v", err)
	}

	if isLeaf == true {
		return true
	}

	return false
}

// directInsert inserts the star with the given ID into the given node inside of the given database
func directInsert(starID int64, db *sql.DB, nodeID int64) {

	// build the query
	query := fmt.Sprintf("UPDATE nodes SET star_id=%d WHERE node_id=%d", starID, nodeID)

	// Execute the query
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("[E] directInsert query: %v", err)
	}
}

func subdivide(db *sql.DB, nodeID int64) {
	log.Printf("[i] Subdividing node %d\n", nodeID)

	boxWidth := getBoxWidth(db, nodeID)
	boxCenter := getBoxCenter(db, nodeID)
	originalDepth := getNodeDepth(db, nodeID)

	log.Printf("[i] original box width: %f", boxWidth)
	log.Printf("[i] original box center: %f", boxCenter)

	// calculate the new positions
	newPosX := boxCenter[0] + (boxWidth / 2)
	newPosY := boxCenter[1] + (boxWidth / 2)
	newNegX := boxCenter[0] - (boxWidth / 2)
	newNegY := boxCenter[1] - (boxWidth / 2)
	newWidth := boxWidth / 2

	log.Printf("[i] new box width: %f", newWidth)
	log.Printf("[i] new box center: [±%f, ±%f]", newPosX, newPosY)

	// create new news with those positions
	newNodeIDA := newNode(db, newPosX, newPosY, newWidth, originalDepth+1)
	newNodeIDB := newNode(db, newPosX, newNegY, newWidth, originalDepth+1)
	newNodeIDC := newNode(db, newNegX, newPosY, newWidth, originalDepth+1)
	newNodeIDD := newNode(db, newNegX, newNegY, newWidth, originalDepth+1)

	// Update the subtrees of the parent node

	// build the query
	query := fmt.Sprintf("UPDATE nodes SET subnode='{%d, %d, %d, %d}', isleaf=FALSE WHERE node_id=%d", newNodeIDA, newNodeIDB, newNodeIDC, newNodeIDD, nodeID)

	// Execute the query
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("[E] subdivide query: %v", err)
	}
}

// getBoxWidth gets the width of the box from the node width the given id
func getBoxWidth(db *sql.DB, nodeID int64) float64 {
	var boxWidth float64

	query := fmt.Sprintf("SELECT box_width FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&boxWidth)
	if err != nil {
		log.Fatalf("[E] getBoxWidth query: %v", err)
	}

	return boxWidth
}

// getBoxWidth gets the center of the box from the node width the given id
func getBoxCenter(db *sql.DB, nodeID int64) []float64 {
	log.Printf("[i] Getting the BoxCenter of node %d\n", nodeID)

	var boxCenter []uint8

	query := fmt.Sprintf("SELECT box_center FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&boxCenter)
	if err != nil {
		log.Fatalf("[E] getBoxCenter query: %v", err)
	}

	boxCenterX, parseErr := strconv.ParseFloat("0", 64)
	boxCenterY, parseErr := strconv.ParseFloat("0", 64)
	if parseErr != nil {
		log.Fatalf("[E] parse boxCenter: %v", err)
		log.Fatalf("[E] parse boxCenter: %s", boxCenter)
	}

	boxCenterFloat := []float64{boxCenterX, boxCenterY}

	return boxCenterFloat
}

// newNode Inserts a new node into the database with the given parameters
func newNode(db *sql.DB, posX float64, posY float64, width float64, depth int64) int64 {

	// build the query creating a new node
	query := fmt.Sprintf("INSERT INTO nodes (box_center, box_width, depth, isleaf) VALUES ('{%f, %f}', %f, %d, TRUE) RETURNING node_id", posX, posY, width, depth)

	var nodeID int64

	// execute the query
	err := db.QueryRow(query).Scan(&nodeID)
	if err != nil {
		log.Fatalf("[E] newNode query: %v", err)
	}

	return nodeID
}

func getBlockingStar(db *sql.DB, nodeID int64) structs.Star2D {
	// 1. get the star id from the node
	// 2. get the stars coordinates from the stars table
	// 3. pack the star and return it

	// get the star id from the node
	var starID int64
	query := fmt.Sprintf("SELECT star_id FROM nodes WHERE node_id=%d", nodeID)
	err := db.QueryRow(query).Scan(&starID)
	if err != nil {
		log.Fatalf("[E] getBlockingStar id query: %v", err)
	}

	fmt.Printf("[i] Getting star with the id %d\n", starID)

	var x, y, vx, vy, m float64

	// get the star from the stars table
	query = fmt.Sprintf("SELECT x, y, vx, vy, m FROM stars WHERE star_id=%d", starID)
	err = db.QueryRow(query).Scan(&x, &y, &vx, &vy, &m)
	if err != nil {
		log.Fatalf("[E] getBlockingStar star query: %v", err)
	}

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

	return star
}

// deleteAll Stars deletes all the rows in the stars table
func deleteAllStars(db *sql.DB) {

	// build the query creating a new node
	query := "DELETE FROM stars WHERE TRUE"

	// execute the query
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("[E] deleteAllStars query: %v", err)
	}
}

// deleteAll Stars deletes all the rows in the nodes table
func deleteAllNodes(db *sql.DB) {

	// build the query creating a new node
	query := "DELETE FROM nodes WHERE TRUE"

	// execute the query
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("[E] deleteAllStars query: %v", err)
	}
}

// getNodeDepth returns the depth of the given node in the tree
func getNodeDepth(db *sql.DB, nodeID int64) int64 {
	log.Printf("[i] Getting the NodeDepth of node %d\n", nodeID)

	// build the query
	query := fmt.Sprintf("SELECT depth FROM nodes WHERE node_id=%d", nodeID)

	var depth int64

	// Execute the query
	err := db.QueryRow(query).Scan(&depth)
	if err != nil {
		log.Fatalf("[E] getNodeDepth query: %v", err)
	}

	return depth
}

// quadrant returns the quadrant into which the given star belongs
func quadrant(db *sql.DB, star structs.Star2D, nodeID int64) int64 {

	// get the center of the node the star is in
	center := getBoxCenter(db, nodeID)
	centerX := center[0]
	centerY := center[1]

	if star.C.X > centerX {
		if star.C.Y > centerY {
			// North East condition
			return 1
		}
		// South East condition
		return 3
	}

	if star.C.Y > centerY {
		// North West condition
		return 0
	}
	// South West condition
	return 2
}

// getQuadrantNodeID returns the id of the requested child-node
// Example: if a parent has four children and quadrant 0 is requested, the function returns the north east child id
func getQuadrantNodeID(db *sql.DB, parentNodeID int64, quadrant int64) int64 {
	log.Printf("[i] Getting the QuadrantNodeId of node %d for the quadrant %d\n", parentNodeID, quadrant)

	//var a string
	var a, b, c, d []uint8

	// get the star from the stars table
	query := fmt.Sprintf("SELECT subnode[1], subnode[2], subnode[3], subnode[4] FROM nodes WHERE node_id=%d", parentNodeID)
	err := db.QueryRow(query).Scan(&a, &b, &c, &d)
	if err != nil {
		log.Fatalf("[E] getQuadrantNodeID star query: %v", err)
	}

	log.Printf("[o] %v", a)

	returnA, _ := strconv.ParseInt(string(a), 10, 64)
	returnB, _ := strconv.ParseInt(string(b), 10, 64)
	returnC, _ := strconv.ParseInt(string(c), 10, 64)
	returnD, _ := strconv.ParseInt(string(d), 10, 64)

	switch quadrant {
	case 0:
		return returnA
	case 1:
		return returnB
	case 2:
		return returnC
	case 3:
		return returnD
	}

	return -1
}

func removeStarFromNode(db *sql.DB, nodeID int64) {

	// build the query
	query := fmt.Sprintf("UPDATE nodes SET star_id=0 WHERE node_id=%d", nodeID)

	// Execute the query
	_, err := db.Query(query)
	if err != nil {
		log.Fatalf("[E] removeStarFromNode query: %v", err)
	}
}

func getStar(db *sql.DB, starID int64) structs.Star2D {
	var x, y, vx, vy, m float64

	// get the star from the stars table
	query := fmt.Sprintf("SELECT x, y, vx, vy, m FROM stars WHERE star_id=%d", starID)
	err := db.QueryRow(query).Scan(&x, &y, &vx, &vy, &m)
	if err != nil {
		log.Fatalf("[E] getStar star query: %v", err)
	}

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

	return star
}
