# optgen

[![GoDoc](https://godoc.org/github.com/lumochift/optgen?status.svg)](https://godoc.org/github.com/lumochift/optgen)
[![Build Status](https://github.com/lumochift/optgen/workflows/Go%20workflow/badge.svg)](https://github.com/lumochift/optgen/actions)
[![codecov](https://codecov.io/gh/lumochift/optgen/branch/master/graph/badge.svg)](https://codecov.io/gh/lumochift/optgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/lumochift/optgen)](https://goreportcard.com/report/github.com/lumochift/optgen)

Optional function generator from struct.

## Limitation

- Flag `-w` currently cannot update generated code, only append code

## Installation

```bash
go install github.com/lumochift/optgen@latest
```

## Usage

- Help function
  
```bash
Usage of ./optgen:
  -all
        generate all fields
  -file string
        path file
  -name string
        struct name
  -tag string
        custom tag (default "opt")
  -w    enable write mode
```

- Generate file with default tags `opt`
  
```go
//thing.go
package foo

type Thing struct {
    Field1 string `opt` 
    Field2 []*int  `opt`
    Field3 map[byte]float64
}
```

For example we have `thing.go` and need to generate functional option for struct `Thing`, only field with `opt` will generated:

```bash
optgen -file thing.go -name Thing
```

Output:

```go
// NewThing returns a *Thing.
func NewThing(opt ...func(*Thing)) *Thing {

        // Prepare a Thing 
        thing := &Thing{}

        // Apply options.
        for _, o := range opt {
                o(thing)
        }

        // Do anything here

        return thing
}


// SetField1 sets the Field1
func SetField1(field1 string) func(*Thing) {
        return func(c *Thing) {
                c.Field1 = field1
        }
}

// SetField2 sets the Field2
func SetField2(field2 []*int) func(*Thing) {
        return func(c *Thing) {
                c.Field2 = field2
        }
}
```
