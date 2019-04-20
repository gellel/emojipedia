package x

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

var re, _ = regexp.Compile(`\(([^()]+)\)`)

var replacements = []string{"[", "", "]", ""}

var replacer = strings.NewReplacer(replacements...)

type Arg struct {
	Argument *Argument
}

type Argument struct {
	Address   uintptr
	Kind      reflect.Kind
	Parameter string
	Pointer   bool
	Position  int
	Name      string
	Slice     bool
	Value     string
	Variadic  bool
}

type Arguments []*Argument

type Deconstruct struct {
	Parameters []string
	Pointer    uintptr
	Name       string
	Type       reflect.Type
	Variadic   bool
}

type Function struct {
	Arguments *Arguments
	Empty     bool
	F         interface{}
	Length    int
	Pointer   uintptr
	Name      string
	Variadic  bool
}

type Functions map[string]*Function

type Runner struct {
	Functions *Functions
}

func NewArg(argument *Argument) (arg *Arg) {
	return &Arg{Argument: argument}
}

func NewArgument(name, value string, position int, pointer uintptr, variadic bool, kind reflect.Kind) (argument *Argument) {
	return &Argument{
		Address:  pointer,
		Kind:     kind,
		Pointer:  strings.Index(value, "*") != -1,
		Position: position,
		Slice:    (kind.String() == "slice"),
		Name:     name,
		Value:    replacer.Replace(value),
		Variadic: variadic}
}

func NewArguments(reflection reflect.Type, pointer uintptr, variadic bool, parameters []string) *Arguments {
	arguments := &Arguments{}
	for i, parameter := range parameters {
		in := reflection.In(i)
		substrings := strings.Split(parameter, " ")
		argument := NewArgument(substrings[0], in.String(), i, pointer, variadic, in.Kind())
		*arguments = append(*arguments, argument)
	}
	return arguments
}

func NewDeconstruct(f interface{}) *Deconstruct {
	reflection := reflect.TypeOf(f)
	pointer := reflect.ValueOf(f).Pointer()
	reference := runtime.FuncForPC(pointer)
	variadic := reflection.IsVariadic()
	name := filepath.Base(reference.Name())
	i := strings.Index(name, ".")
	for i > -1 {
		name = name[(i + 1):]
		i = strings.Index(name, ".")
	}
	parameters := NewParameters(reference.FileLine(pointer))
	return &Deconstruct{
		Parameters: parameters,
		Pointer:    pointer,
		Name:       name,
		Type:       reflection,
		Variadic:   variadic}
}

func NewFunction(f interface{}) (function *Function) {
	deconstruct := NewDeconstruct(f)
	arguments := NewArguments(deconstruct.Type, deconstruct.Pointer, deconstruct.Variadic, deconstruct.Parameters)
	length := len(*arguments)
	return &Function{
		Arguments: arguments,
		Empty:     length == 0,
		F:         f,
		Length:    length,
		Pointer:   deconstruct.Pointer,
		Name:      deconstruct.Name,
		Variadic:  deconstruct.Variadic}
}

func NewFunctions(f ...interface{}) (functions *Functions) {
	functions = &Functions{}
	for _, x := range f {
		n := NewFunction(x)
		(*functions)[strings.ToUpper(n.Name)] = n
	}
	return functions
}

func NewParameters(file string, line int) (arguments []string) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	contents := bytes.Split(b, []byte("\n"))
	substring := string(contents[line-1])
	matches := re.FindAllStringSubmatch(substring, 1)
	if len(matches) != 0 {
		arguments = strings.Split(matches[0][1], ",")
	}
	return arguments
}

func NewRunner(f ...interface{}) (runner *Runner) {
	return &Runner{
		Functions: NewFunctions(f...)}
}

func (arg *Arg) Is(key string) (ok bool) {
	ok = arg.Argument != nil && arg.Argument.Is(key)
	return ok
}

func (argument *Argument) Is(key string) (ok bool) {
	fmt.Println(argument.Value)
	ok = strings.ToUpper(argument.Value) == strings.ToUpper(key)
	return ok
}

func (argument *Argument) IsEach(values ...string) {}

func (arguments *Arguments) Bounds(i int) (ok bool) {
	ok = ((i > -1) && (i < len(*arguments)))
	return ok
}

func (arguments *Arguments) Each(function func(i int, argument *Argument)) *Arguments {
	for i, argument := range *arguments {
		function(i, argument)
	}
	return arguments
}

func (arguments *Arguments) Get(i int) (argument *Argument, ok bool) {
	if ok = arguments.Bounds(i); ok {
		argument = (*arguments)[i]
	}
	return argument, ok
}

func (arguments *Arguments) Length() (length int) {
	length = len(*arguments)
	return length
}

func (arguments *Arguments) Peek(i int) (arg *Arg) {
	arg = &Arg{}
	if argument, ok := arguments.Get(i); ok {
		arg.Argument = argument
	}
	return arg
}

func (function *Function) Set(f interface{}) *Function {
	*function = *NewFunction(f)
	return function
}

func (functions *Functions) Contains(function *Function) (ok bool) {
	ok = functions.Has(function.Name)
	return ok
}

func (functions *Functions) Fetch(key string) (function *Function) {
	function, _ = functions.Get(key)
	return function
}

func (functions *Functions) Get(key string) (function *Function, ok bool) {
	function, ok = (*functions)[strings.ToUpper(key)]
	return function, ok
}

func (functions *Functions) Has(key string) (ok bool) {
	_, ok = (*functions)[strings.ToUpper(key)]
	return ok
}

func (functions *Functions) Set(f ...interface{}) *Functions {
	*functions = *NewFunctions(f...)
	return functions
}

func (runner *Runner) Get(key string) (f func(i ...interface{}), ok bool) {
	function, ok := runner.Functions.Get(key)
	if ok && function.Variadic {
		if argument, ok := function.Arguments.Get(0); ok {
			if ok = argument.Is("interface"); ok {
				f = function.F.(func(i ...interface{}))
			}
		}
	}
	return f, ok
}

func (runner *Runner) Set(f ...interface{}) *Runner {
	*runner = *NewRunner(f...)
	return runner
}
