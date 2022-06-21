package gograylog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	grayLogDateFormat string = "2006-01-02T15:04:05.000Z"
)

type Query struct {
	Host      string
	Query     string
	Streamid  string
	Fields    []string
	Limit     int
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
	if len(q.Fields) > 0 {
		data["fields_in_order"] = q.Fields
	}
	data["limit"] = q.Limit

	return data
}

func (q Query) ToURL(from, to time.Time) string {
	params := url.Values{}

	params.Add("q", q.Query)
	params.Add("fields_in_order", strings.Join(q.Fields, ", "))

	params.Add("timerange", "absolute")
	params.Add("from", from.Format(grayLogDateFormat))
	params.Add("to", to.Format(grayLogDateFormat))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.Host, q.Streamid, params.Encode())
}
