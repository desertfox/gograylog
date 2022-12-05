package gograylog

import (
	"fmt"
	"testing"
)

func Test_buildBodyData(t *testing.T) {

	var q query = query{
		host:      "https://desertfox.dev",
		query:     "error",
		streamid:  "somehash",
		frequency: 15,
	}

	t.Log(q.url())
	t.Log(q.request())
	t.Log(fmt.Sprintf("%s", q.data()))

}
