package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/boltdb/bolt"
	"gopkg.in/yaml.v2"
)

type mapping struct {
	Path string `yaml:"path" json:"path"`
	Url  string `yaml:"url" json:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if redirect, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, redirect, http.StatusMovedPermanently)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(mappingsYaml []byte, fallback http.Handler) (http.HandlerFunc, error) {

	mappings, err := parseYaml(mappingsYaml)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(mappings)

	return MapHandler(pathsToUrls, fallback), nil
}

func buildMap(mappings []mapping) map[string]string {
	pathsToUrls := make(map[string]string)
	for _, mapper := range mappings {
		pathsToUrls[mapper.Path] = mapper.Url
	}
	return pathsToUrls
}

func parseYaml(yml []byte) ([]mapping, error) {
	var mappings []mapping
	err := yaml.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML document: %v", err)
	}
	return mappings, nil
}

func parseJson(yml []byte) ([]mapping, error) {
	var mappings []mapping
	err := json.Unmarshal(yml, &mappings)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON document: %v", err)
	}
	return mappings, nil
}

func JSONHandler(mappingsJson []byte, fallback http.Handler) (http.HandlerFunc, error) {

	var mappings []mapping
	err := json.Unmarshal(mappingsJson, &mappings)
	if err != nil {
		return nil, err
	}
	pathsToUrls := buildMap(mappings)

	return MapHandler(pathsToUrls, fallback), nil
}

func BoltHandler(databaseFile string, fallback http.Handler) (http.HandlerFunc, error) {

	return func(w http.ResponseWriter, r *http.Request) {

		db, err := bolt.Open(databaseFile, 0600, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		defer db.Close()

		var redirect string
		err = db.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("mappings"))
			redirectFromDB := bucket.Get([]byte(r.URL.Path))
			redirect = string(redirectFromDB)
			return nil
		})

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

		if len(redirect) != 0 {
			http.Redirect(w, r, redirect, http.StatusMovedPermanently)
			return
		} else {
			fallback.ServeHTTP(w, r)
		}
	}, nil

}
