```go run
package main

import "fmt"

func counter() func() int {
	n := 0
	return func() int {
		n++
		return n
	}
}

func main() {
	c := counter()
	c()
	fmt.Println(c())
}
```

```output
2
```
