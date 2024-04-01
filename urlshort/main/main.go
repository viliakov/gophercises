package main

import (
	"fmt"
	"net/http"

	"github.com/viliakov/gophercises/urlshort"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	json := `
[
	{"path": "/json-godoc", "url": "https://pkg.go.dev/encoding/json"},
	{"path": "/json", "url": "https://www.json.org/json-en.html"}
]
`
	jsonHandler, err := urlshort.JSONHandler([]byte(json), yamlHandler)
	if err != nil {
		panic(err)
	}

	databaseFile := "urlShortener.db"
	boltHandler, err := urlshort.BoltHandler(databaseFile, jsonHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", boltHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
