package keywords

import (
	"github.com/gellel/emojipedia/lexicon"
	"github.com/gellel/emojipedia/slice"
)

// New instantiates a new empty Keywords pointer.
func New() *Keywords {
	return &Keywords{&lexicon.Lexicon{}}
}

type keywords interface {
	Add(key string, names ...string) *Keywords
	Each(f func(slice *slice.Slice)) *Keywords
	Fetch(key string) *slice.Slice
	Get(key string) (*slice.Slice, bool)
	Has(key string) bool
	Keys() *slice.Slice
	Len() int
	Remove(key string) bool
	Values() *slice.Slice
}

// Keywords is a map-like struct with methods used to perform traversal and retrieval of slice.Slice pointers.
type Keywords struct {
	lexicon *lexicon.Lexicon
}

// Add method adds one or more strings to the struct using the key reference to update or create the associated slice.
func (pointer *Keywords) Add(key string, names ...string) *Keywords {
	if pointer.lexicon.Has(key) == false {
		pointer.lexicon.Add(key, slice.New())
	}
	s := pointer.lexicon.Fetch(key).(*slice.Slice)
	for _, name := range names {
		s.Append(name)
	}
	return pointer
}

// Each method executes a provided function once for each slice.Slice pointer.
func (pointer *Keywords) Each(f func(key string, slice *slice.Slice)) *Keywords {
	pointer.lexicon.Each(func(key string, i interface{}) {
		f(key, i.(*slice.Slice))
	})
	return pointer
}

// Fetch retrieves the slice.Slice pointer held by the argument key. Panics if key does not exist.
func (pointer *Keywords) Fetch(key string) *slice.Slice {
	property, _ := pointer.Get(key)
	return property
}

// Get returns the slice.Slice pointer held by the argument key and a boolean indicating if it was successfully retrieved.
// Panics if cannot convert to slice.Slice pointer.
func (pointer *Keywords) Get(key string) (*slice.Slice, bool) {
	property, ok := pointer.lexicon.Get(key)
	return property.(*slice.Slice), ok
}

// Has method checks that a given key exists in the Keywords.
func (pointer *Keywords) Has(key string) bool {
	return pointer.lexicon.Has(key)
}

// Keys method returns a slice.Slice of a given Keywords' own property names, in the same order as we get with a normal loop.
func (pointer *Keywords) Keys() *slice.Slice {
	slice := slice.New()
	pointer.lexicon.Each(func(key string, i interface{}) {
		slice.Append(key)
	})
	return slice
}

// Len method returns the number of elements in the Keywords.
func (pointer *Keywords) Len() int {
	return pointer.lexicon.Len()
}

// Remove method removes a entry from the Keywords if it exists. Returns a boolean to confirm if it succeeded.
func (pointer *Keywords) Remove(key string) bool {
	return pointer.lexicon.Remove(key)
}

// Values method returns a Slice of a given Keywords's own enumerable property values,
// in the same order as that provided by a for...in loop.
func (pointer *Keywords) Values() *slice.Slice {
	slice := slice.New()
	pointer.lexicon.Each(func(key string, i interface{}) {
		slice.Append(i)
	})
	return slice
}