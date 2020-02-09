# optgen

Optional function generator from struct.

## Installation

```
go get github.com/lucmohift/optgen
```

## How to use

```go
//thing.go
package foo

type Thing struct {
    Field1 string `opt` 
    Field2 []*int  `opt`
    Field3 map[byte]float64
}

type ThingThong struct {
    F1 string `opt` 
    F2 []int  `opt`
    F3 map[byte]float64
}

type Data string
```

For example we have `thing.go` and need to generate functional option for struct `Thing`:

```bash
optgen -file thing.go -name Thing
```

Output:

```go
// Option is a Thing configurator to be supplied to NewThing() function.
type Option func(*Thing)


// NewThing returns a new Thing.
func NewThing(options ...Option) (*Thing, error) {

        // Prepare a Thing with default host.
        thing := &Thing{}

        // Apply options.
        for _, option := range options {
                option(thing)
        }

        // Do anything here

        return thing, nil
}


// SetField1 sets the Field1
func SetField1(field1 string) Option {
        return func(c *Thing) {
                c.Field1 = field1
        }
}

// SetField2 sets the Field2
func SetField2(field2 []*int) Option {
        return func(c *Thing) {
                c.Field2 = field2
        }
}
```
