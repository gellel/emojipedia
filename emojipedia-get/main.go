package get

import (
	"strings"

	"github.com/gellel/emojipedia/manifest"

	categories "github.com/gellel/emojipedia/emojipedia-get/get-categories"
	emoji "github.com/gellel/emojipedia/emojipedia-get/get-emoji"
	keywords "github.com/gellel/emojipedia/emojipedia-get/get-keywords"
	subcategories "github.com/gellel/emojipedia/emojipedia-get/get-subcategories"
)

var Export = get

var programs = map[string](func(m *manifest.Manifest, previous, options []string)){
	categories.Key:    categories.Main,
	emoji.Key:         emoji.Main,
	keywords.Key:      keywords.Main,
	subcategories.Key: subcategories.Main}

func main(m *manifest.Manifest, previous []string, options []string) {
	var argument string
	if len(options) != 0 {
		argument = strings.ToUpper(options[0])
	}
	if f, ok := programs[argument]; ok {
		f(m, append(previous, argument), options[1:])
	}
}

func get(arguments ...interface{}) {
	main(arguments[0].(*manifest.Manifest), arguments[1].([]string), arguments[2].([]string))
}
