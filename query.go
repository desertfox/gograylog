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

type query struct {
	host      string
	query     string
	streamid  string
	fields    []string
	limit     int
	frequency int
}

func (q query) url() string {
	return fmt.Sprintf("%s/api/views/search/messages", q.host)
}

func (q query) data() io.Reader {
	data, err := json.Marshal(q.request())
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(data)
}

func (q query) request() map[string]interface{} {
	var data map[string]interface{} = make(map[string]interface{}, 3)

	data["streams"] = []string{q.streamid}
	data["timerange"] = map[string]string{
		"type":  "relative",
		"range": strconv.Itoa(q.frequency * 60),
	}
	data["query_string"] = map[string]string{
		"type":         "elasticsearch",
		"query_string": q.query,
	}
	if len(q.fields) > 0 {
		data["fields_in_order"] = q.fields
	}
	data["limit"] = q.limit

	return data
}

func (q query) string(from, to time.Time) string {
	params := url.Values{}

	params.Add("q", q.query)
	params.Add("fields_in_order", strings.Join(q.fields, ", "))

	params.Add("timerange", "absolute")
	params.Add("from", from.Format(grayLogDateFormat))
	params.Add("to", to.Format(grayLogDateFormat))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.host, q.streamid, params.Encode())
}
