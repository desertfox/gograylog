package gograylog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

type Query struct {
	Host      string
	Query     string
	Streamid  string
	Frequency int
}

func (q Query) URL() string {
	return fmt.Sprintf("%s/api/views/search/messages", q.Host)
}

func (q Query) BodyData() io.Reader {
	data, err := json.Marshal(q.buildBodyData())
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

func (q Query) buildBodyData() map[string]interface{} {
	var streams []string = []string{q.Streamid}

	var relative map[string]string = make(map[string]string)
	relative["type"] = "relative"
	relative["range"] = strconv.Itoa(q.Frequency * 60)

	var queryString map[string]string = make(map[string]string)
	queryString["type"] = "elasticsearch"
	queryString["query_string"] = q.Query

	var data map[string]interface{} = make(map[string]interface{}, 3)
	data["streams"] = streams
	data["timerange"] = relative
	data["query_string"] = queryString

	return data
}
