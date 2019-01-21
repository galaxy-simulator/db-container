package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

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
