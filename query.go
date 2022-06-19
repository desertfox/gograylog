package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

var (
//grayLogDateFormat   string = "2006-01-02T15:04:05.000Z"
//relativeStrTempalte string = "%v/api/search/universal/relative?%v"
//absoluteStrTempalte string = "%v/api/search/universal/absolute?%v"

)

type Query struct {
	Host      string
	Query     string
	Streamid  string
	Frequency int
}

func (q Query) URL() string {
	return fmt.Sprintf("%s/views/search/messages", q.Host)
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

/*

/views/search/messages

 "timerange": {
    "type": "relative",
    "range": 30000
  }

{
	"streams": [
	  "000000000000000000000001"
	],
	"timerange": [
	  "absolute",
	  {
		"from": "2020-12-01T00:00:00.000Z",
		"to": "2020-12-01T15:00:00.000Z"
	  }
	],
	"query_string": { "type":"elasticsearch", "query_string":"your_query" }
  }

*/
