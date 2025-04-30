package urlshort

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/yaml.v2"
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
		foundUrl, isFound := mapHandlerSearchURLS(pathsToUrls, r)
		switch isFound {
		case true:
			http.Redirect(w, r, foundUrl, http.StatusSeeOther)
		default:
			fallback.ServeHTTP(w, r)
		}
	}
}

func mapHandlerSearchURLS(pathsToUrls map[string]string, r *http.Request) (foundUrl string, isFound bool) {
	isFound = false
	foundUrl = ""
	pathStrings := getPathsFromRequest(r)
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

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsed, err := ParseYAML(yml)

	if err != nil {
		return nil, err
	}

	if len(parsed) == 0 {
		return nil, errors.New("yaml returned no information")
	}
	for _, parsed := range parsed {
		fmt.Println(parsed.repr())
	}
	return func(w http.ResponseWriter, r *http.Request) {
		foundUrl, isFound := yamlHandlerSearchUrls(parsed, r)
		switch isFound {
		case true:
			http.Redirect(w, r, foundUrl, http.StatusSeeOther)
			return
		default:
			fallback.ServeHTTP(w, r)
			return
		}
	}, nil
}

func yamlHandlerSearchUrls(pathToUrls []*PathToUrl, r *http.Request) (foundUrl string, isFound bool) {
	isFound = false
	foundUrl = ""
	pathStrings := getPathsFromRequest(r)
	for _, pathString := range pathStrings {
		fmt.Println("evaluating path: ", pathString)
		formedPath := fmt.Sprintf("/%s", pathString)
		fmt.Println("current formedPath: ", formedPath)
		for _, pathToUrl := range pathToUrls {
			fmt.Println("current pathToURlPath", pathToUrl.Path)
			if formedPath == pathToUrl.Path {
				foundUrl = pathToUrl.URL
				isFound = true
				return
			}
		}
	}
	return
}

func ParseYAML(data []byte) (parsed []*PathToUrl, err error) {
	parsedData := []PathToUrl{}
	err = yaml.Unmarshal(data, &parsedData)
	if err != nil {
		return nil, err
	}
	for _, record := range parsedData {
		parsed = append(parsed, &record)
	}
	return parsed, nil
}

type PathToUrl struct {
	Path string `yaml:"path"`
	URL  string `yaml:"url"`
}

func (p *PathToUrl) repr() string {
	return fmt.Sprintf("Path: %s\nURL %s", p.Path, p.URL)
}
func getPathsFromRequest(r *http.Request) []string {
	url := r.URL
	path := url.Path
	pathStrings := strings.Split(path, "/")
	return pathStrings
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
