# db-container

This Repo contains the main "database" running an http-api exposing the quadtree.

## API-Endpoints


## Files

Simple description of each individual function

### export.go

```
func export(treeindex int64) error {
    // exports the tree with the given tree index by writing it's json 
    // to /exports/{treeindex}.json
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
    // handles request to /export/{treeindex} and uses the export function
    // to export the selected tree
}
```

### update.go

```
func updateCenterOfMassHandler(w http.ResponseWriter, r *http.Request) {
    // updates the center of mass of every node in the tree width the given index
}

func updateTotalMassHandler(w http.ResponseWriter, r *http.Request) {
    // updates the total mass of every node in the tree with the given index
}
```

### import.go

```
func fastInsertHandler(w http.ResponseWriter, r *http.Request) {
    // creates a new tree by reading from /db/{filename}.json and inserting the tree
    // from that file into the treeArray
    // The function returns the index in which it inserted the tree
}
```

### manage.go

```
func newTreeHandler(w http.ResponseWriter, r *http.Request) {
    // creates a new empty tree with the given width
}

func newTree(width float64) []byte {
    // creates a new tree and returns it json formatted as a bytestring
}

func printAllHandler(w http.ResponseWriter, r *http.Request) { 
    // prints all the trees in the treeArray in json format
}

func insertStarHandler(w http.ResponseWriter, r *http.Request) {
    // inserts the given star into the tree
}

func starlistHandler(w http.ResponseWriter, r *http.Request) {
    // lists all the stars in the given tree
}

func dumptreeHandler(w http.ResponseWriter, r *http.Request) {
    // returns the tree with the given index in json format 
}
```

### metrics.go

```
func metricHandler(w http.ResponseWriter, r *http.Request) {
    // locally published the databases metrics
}

func pushMetricsNumOfStars(host string, treeindex int64) {
    // pushes the metrics of the tree with the given index to the metric-bundler
    // using the given host and 
}
```

### update.go

```
func updateCenterOfMassHandler(w http.ResponseWriter, r *http.Request) {
    // updates the center of mass of every node in the tree with the given index
}

func updateTotalMassHandler(w http.ResponseWriter, r *http.Request) {
    // updates the total mass in the tree with the given index
}
```