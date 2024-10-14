package tool

import (
	"fmt"
)

func ExampleExtractFormat() {
	fmt.Println(ExtractFormat("hello/world.txt"))
	fmt.Println(ExtractFormat("hello/world.txt#json"))
	// Output:
	// hello/world.txt txt
	// hello/world.txt json
}
