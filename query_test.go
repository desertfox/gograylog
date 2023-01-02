package gograylog

import (
	"fmt"
)

func ExampleQuery() {
	q := Query{
		Host:        "https://desertfox.dev",
		QueryString: "error",
		StreamID:    "somehash",
		Frequency:   15,
	}

	fmt.Println(q, q.endpoint())
	//Output: {https://desertfox.dev error somehash [] 0 15} https://desertfox.dev/api/system/sessions
}
