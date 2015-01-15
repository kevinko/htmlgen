# Simple HTML generation for Go

This allows the generation of HTML programmatically.

See the following example:

```
package main

import (
	H "github.com/kevinko/htmlgen"
	"os"
)

func main() {
	root := H.NewRoot()
	body := root.Body()
	body.H1().T(`Hello, world!`)
	body.P().TV(`How are you, $firstName $lastName?`)

	env := H.Environment{
		"firstName": H.StringValue("John"),
		"lastName":  H.StringValue("Doe"),
	}

	H.Write(os.Stdout, root, env)
}
```

This results in a binary that generates the output:

```
<!DOCTYPE html><html><body><h1>Hello, world!</h1><p>How are you, John Doe?</p></body></html>
```

Using a 1000x10 table generation benchmark, performance is similar to
Mako on Python.
