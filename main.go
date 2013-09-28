package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var port = flag.Int("port", 8000, "Port to serve the sprinkler control interface")
var pin = flag.String("pin_path", "/sys/class/gpio/gpio48/value",
	"Path to pin's value file. It's expected to have direction set to 'out'")

func indexHandler(w http.ResponseWriter, r *http.Request) {
	state, err := ioutil.ReadFile(*pin)
	if err != nil {
		fmt.Fprintf(w, "Error writing to [%q]: %q",
			html.EscapeString(*pin),
			html.EscapeString(err.Error()))
		return
	}
	fmt.Fprintf(w, "Current state is: %q", html.EscapeString(strings.TrimSpace(string(state))))
}

func switchHandler(w http.ResponseWriter, r *http.Request) {
	stateStr := r.FormValue("state")
	stateStr = strings.TrimSpace(strings.ToLower(stateStr))
	state := stateStr == "on" || stateStr == "1" || stateStr == "true"
	val := "0"
	if state {
		val = "1"
	}
	if err := ioutil.WriteFile(*pin, []byte(val), 0644); err != nil {
		fmt.Fprintf(w, "Error writing to [%q]: %q",
			html.EscapeString(*pin),
			html.EscapeString(err.Error()))
		return
	}
	fmt.Fprintf(w, "OK!\n")
}

func main() {
	flag.Parse()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/switch", switchHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
