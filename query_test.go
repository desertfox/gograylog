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

	fmt.Println(q, q.endpoint("https://desertfox.dev"))
	//Output: {error somehash [] 0 15} https://desertfox.dev/api/views/search/messages
}
