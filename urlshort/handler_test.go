package urlshort

import (
	"fmt"
	"testing"
)

func TestParseYAML(t *testing.T) {

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	data, err := ParseYAML([]byte(yaml))
	if err != nil {
		t.Errorf("didn't expect error got err: %v", err)
	}

	for _, row := range data {
		url := row.URL
		path := row.Path

		fmt.Println("url: ", url)
		fmt.Println("path: ", path)
	}

}
