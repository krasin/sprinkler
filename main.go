package main

import (
	"flag"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var port = flag.Int("port", 8000, "Port to serve the sprinkler control interface")
var pin = flag.String("pin_path", "/sys/class/gpio/gpio48/value",
	"Path to pin's value file. It's expected to have direction set to 'out'")

type IndexData struct {
	State string
}

const templ = `<html>
<head>
<script>
function redirect(url) {
  window.location = url;
}
</script>
</head>
<body>
Current state is {{.State}}.
<br/>
<br/>
<button name="on_btn" onclick="redirect('/switch?state=on')">On</button>
<button name="off_btn" onclick="redirect('/switch?state=off')">Off</button>
</body>
</html>
`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(*pin)
	if err != nil {
		fmt.Fprintf(w, "Error writing to [%q]: %q",
			html.EscapeString(*pin),
			html.EscapeString(err.Error()))
		return
	}
	state := strings.TrimSpace(string(data))

	t := template.New("Index template")
	if t, err = t.Parse(templ); err != nil {
		log.Fatal("templ parse failed: ", err)
	}
	err = t.Execute(w, &IndexData{
		State: state,
	})
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
	http.Redirect(w, r, "/", 307)
}

func main() {
	flag.Parse()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/switch", switchHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
