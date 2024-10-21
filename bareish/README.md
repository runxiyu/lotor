# go-barish

An implementation of a [BARE](https://baremessages.org)-like message format for
Go. This is a fork of [go-bare](https://git.sr.ht/~sircmpwn/go-bare).

**Status**

This mostly works, but you may run into some edge cases with union types.

## Changes from upstream

Not much yet.

We intend to remove all external dependencies and remove some limitations
around integer sizes.

## Code generation

An example is provided in the `examples` directory. Here is a basic
introduction:

```
$ cat schema.bare
type Address {
	address: [4]string
	city: string
	state: string
	country: string
}
$ go run git.sr.ht/~runxiyu/go-bareish/cmd/gen -p models schema.bare models/gen.go
```

Then you can write something like the following:

```go
import "models"

/* ... */

bytes := []byte{ /* ... */ }
var addr Address
err := addr.Decode(bytes)
```

You can also add custom types and skip generating them by passing the `-s
TypeName` flag to gen, then providing your own implementation. For example, to
rig up time.Time with a custom "Time" BARE type, add this to your BARE schema:

```
type Time string
```

Then pass `-s Time` to gen, and provide your own implementation of Time in the
same package. See `examples/time.go` for an example of such an implementation.

## Marshal usage

For many use-cases, it may be more convenient to write your types manually and
use Marshal and Unmarshal directly. If you choose this approach, you may also
use `git.sr.ht/~runxiyu/go-bareish/schema.SchemaFor` to generate a BARE schema
language document describing your structs.

```go
package main

import (
    "fmt"
    "git.sr.ht/~runxiyu/go-bareish"
)

// type Coordinates {
//    x: int
//    y: int
//    z: int
//    q: optional<int>
// }
type Coordinates struct {
    X uint
    Y uint
    Z uint
    Q *uint
}

func main() {
    var coords Coordinates
    payload := []byte{0x01, 0x02, 0x03, 0x01, 0x04}
    err := bare.Unmarshal(payload, &coords)
    if err != nil {
        panic(err)
    }
    fmt.Printf("coords: %d, %d, %d (%d)\n",
        coords.X, coords.Y, coords.Z, *coords.Q) /* coords: 1, 2, 3 (4) */
}
```

### Unions

To use union types, you need to define an interface to represent the union of
possible values, and this interface needs to implement `bare.Union`:

```go
type Person interface {
	Union
}
```

Then, for each possible union type, implement the interface:

```go
type Employee struct { /* ... */ }
func (e Employee) IsUnion() {}

type Customer struct { /* ... */ }
func (c Customer) IsUnion() {}
```

The IsUnion function is necessary to make the type compatible with the Union
interface. Then, to marshal and unmarshal using this union type, you need to
tell go-bare about your union:

```go
func init() {
    // The first argument is a pointer of the union interface, and the
    // subsequent arguments are values of each possible subtype, in ascending
    // order of union tag:
    bare.RegisterUnion((*Person)(nil)).
      Member(*new(Employee), 0).
      Member(*new(Customer), 1)
}
```

This is all done for you if you use code generation.
