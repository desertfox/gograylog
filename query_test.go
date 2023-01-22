package gograylog

import (
	"fmt"
)

func ExampleQuery() {
	q := Query{
		QueryString: "error",
		StreamID:    "somehash",
		Frequency:   15,
	}

	fmt.Println(q)
	//Output: {error somehash [] 0 15}
}
