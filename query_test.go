package gograylog

import (
	"fmt"
	"testing"
)

func Test_buildBodyData(t *testing.T) {

	var q Query = Query{
		Host:      "https://desertfox.dev",
		Query:     "error",
		Streamid:  "somehash",
		Frequency: 15,
	}

	t.Log(q.URL())
	t.Log(q.buildBodyData())
	t.Log(fmt.Sprintf("%s", q.BodyData()))

}
