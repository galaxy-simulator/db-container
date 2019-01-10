package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/marcusolsson/tui-go"

	"git.darknebu.la/GalaxySimulator/db-container/shell/structs"
)

func main() {
	parameters := structs.NewEnv("localhost:8042", "data.csv", 10, 0)
	boundary := structs.NewBoundary(0, 0, 4000)

	// define the history box displayed in the top part of the window
	history := tui.NewVBox()
	historyScroll := tui.NewScrollArea(history)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)

	// define the input bar displayed on the bottom window edge
	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	// define a root container containing all the containers and maximizing them in the given space
	root := tui.NewVBox(historyBox, inputBox)
	root.SetSizePolicy(tui.Expanding, tui.Expanding)

	// on submission of the input box, add the text to the history box
	input.OnSubmit(func(e *tui.Entry) {

		// add the input to the history
		history.Append(tui.NewHBox(tui.NewLabel(e.Text())))
		handleInput(history, e.Text(), parameters, boundary)
		input.SetText("")
	})

	ui, err := tui.New(root)
	if err != nil {
		panic(err)
	}

	ui.SetKeybinding("Esc", func() { ui.Quit() })

	if err := ui.Run(); err != nil {
		panic(err)
	}
}

func faultyInput(history *tui.Box) {
	// print an error message if the input cannot be assigned to an action
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-30s", "", "Faulty Input!")),
		),
	)
}

func printEnvironment(history *tui.Box, parameters *structs.Env) {
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30s", "", "url", parameters.Url())),
		),
	)
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30s", "", "data", parameters.Data())),
		),
	)
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30d", "", "amount", parameters.Amount())),
		),
	)
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30d", "", "treeindex", parameters.Treeindex())),
		),
	)
}

func printBoundary(history *tui.Box, boundary *structs.Boundary) {
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30d", "", "x", boundary.X())),
		),
	)
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30d", "", "y", boundary.Y())),
		),
	)
	history.Append(
		tui.NewHBox(
			tui.NewLabel(fmt.Sprintf("%10s%-15s%-30d", "", "width", boundary.Width())),
		),
	)
}

func setEnvironment(history *tui.Box, parameters *structs.Env, key string, value string) {
	switch key {
	case "url":
		parameters.SetUrl(value)
		break
	case "data":
		parameters.SetData(value)
		break
	case "amount":
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		parameters.SetAmount(valueInt)
		break
	case "treeindex":
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		parameters.SetTreeindex(valueInt)
		break
	default:
		faultyInput(history)
	}
}

func setBoundary(history *tui.Box, boundary *structs.Boundary, key string, value string) {
	switch key {
	case "x":
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		boundary.SetX(valueInt)
		break
	case "y":
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		boundary.SetY(valueInt)
		break
	case "width":
		valueInt, _ := strconv.ParseInt(value, 10, 64)
		boundary.SetWidth(valueInt)
		break
	default:
		faultyInput(history)
	}
}

func handleInput(history *tui.Box, input string, parameters *structs.Env, boundary *structs.Boundary) {
	if input == "help" {
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-35s%-30s", "", "print <env|boundary|all>", "print the according local values")),
			),
		)
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-35s%-30s", "", "set <env|boundary> [KEY] [VALUE]", "set the local values")),
			),
		)
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-35s%-30s", "", "new", "Create a new tree using the values defined in env and boundary")),
			),
		)
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-35s%-30s", "", "<ESC>", "quit")),
			),
		)
		return
	}

	if strings.HasPrefix(input, "print") {
		key := ""

		if len(input) > 6 {
			key = strings.Split(input[6:], " ")[0]
		} else {
			history.Append(
				tui.NewHBox(
					tui.NewLabel(fmt.Sprintf("%10s%-35s", "", "[env|bounadry|all]")),
				),
			)
			return
		}

		switch key {
		case "env":
			printEnvironment(history, parameters)
			break
		case "boundary":
			printBoundary(history, boundary)
			break
		case "all":
			printEnvironment(history, parameters)
			printBoundary(history, boundary)
			break
		default:
		}
		return
	}

	// test if the input string starts with 'set'
	if strings.HasPrefix(input, "set") {
		// trim of the "set" off the string
		parameter := strings.Split(input[4:], " ")[0]
		key := strings.Split(input[4:], " ")[1]
		value := strings.Split(input[4:], " ")[2]

		switch parameter {
		case "env":
			setEnvironment(history, parameters, key, value)
			break
		case "boundary":
			setBoundary(history, boundary, key, value)
			break
		}

		return
	}

	if input == "new" {
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-30s", "", "Creating a new tree")),
			),
		)

		// Generate the request url
		requestUrl := fmt.Sprintf("http://%s/new", parameters.Url())

		// Bundle the post request data
		data := []byte(fmt.Sprintf("x=%d&y=%d&w=%d", boundary.X(), boundary.Y(), boundary.Width()))

		// make the new request the the endpoint defined above using the data bundled above
		req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(data))
		if err != nil {
			log.Fatal("Error reading request.", err)
		}

		// Send the request
		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error reading response.", err)
		}

		// return the http response code
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-30s", "", resp.Status)),
			),
		)
		return
	}

	if input == "insert" {
		infoMessage := fmt.Sprintf("Inserting %d stars from %s into tree Nr. %d defined in %s.", parameters.Amount(), parameters.Data(), parameters.Treeindex(), parameters.Url())
		history.Append(
			tui.NewHBox(
				tui.NewLabel(fmt.Sprintf("%10s%-30s", "", infoMessage)),
			),
		)

		// open parameters.Data()
		dat, err := ioutil.ReadFile(parameters.Data())
		if err != nil {
			log.Fatal("Error reading file")
		}

		// parse the data using a csv reader
		csvData := csv.NewReader(strings.NewReader(string(dat)))

		for i := 0; i < int(parameters.Amount()); i++ {
			record, err := csvData.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("Error Reading the csv data: ", err)
			}

			history.Append(
				tui.NewHBox(
					tui.NewLabel(fmt.Sprintf("%10s%-30s", "", record)),
				),
			)

			// Generate the request url
			requestUrl := fmt.Sprintf("http://%s/insert/%d", parameters.Url(), parameters.Treeindex())

			history.Append(
				tui.NewHBox(
					tui.NewLabel(fmt.Sprintf("%10s%-30s", "requesturl: ", requestUrl)),
				),
			)

			x, _ := strconv.ParseFloat(record[0], 64)
			y, _ := strconv.ParseFloat(record[0], 64)

			// Bundle the post request data
			data := []byte(fmt.Sprintf("x=%f&y=%f&vx=0&vy=0&m=10", x, y))

			// make the new request the the endpoint defined above using the data bundled above
			req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(data))
			if err != nil {
				log.Fatal("Error reading request.", err)
			}

			// Send the request
			client := http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal("Error reading response.", err)
			}

			// return the http response code
			history.Append(
				tui.NewHBox(
					tui.NewLabel(fmt.Sprintf("%10s%-30s", "", resp.Status)),
				),
			)
		}

		return
	}

	faultyInput(history)
}
