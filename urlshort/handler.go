package urlshort

import (
	"fmt"
	"net/http"
	"strings"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	//	TODO: Implement this...
	return func(w http.ResponseWriter, r *http.Request) {
		foundUrl, isFound := searchUrls(pathsToUrls, r)
		switch isFound {
		case true:
			http.Redirect(w, r, foundUrl, http.StatusSeeOther)
		default:
			fallback.ServeHTTP(w, r)
		}
	}
}
func searchUrls(pathsToUrls map[string]string, r *http.Request) (foundUrl string, isFound bool) {
	isFound = false
	foundUrl = ""
	url := r.URL
	path := url.Path
	pathStrings := strings.Split(path, "/")
	for _, pathString := range pathStrings {
		formedString := fmt.Sprintf("/%s", pathString)
		if _, ok := pathsToUrls[formedString]; ok {
			foundUrl = pathsToUrls[formedString]
		}
	}
	switch foundUrl != "" {
	case true:
		isFound = true
		return
	default:
		return
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
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// TODO: Implement this...
	return nil, nil
}
