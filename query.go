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

//date format required by graylogs
const (
	GrayLogDateFormat string = "2006-01-02T15:04:05.000Z"
	MessagesPath      string = "api/system/sessions"
)

type query struct {
	host, query, streamid string
	fields                []string
	limit, frequency      int
}

func (q query) url() string {
	return fmt.Sprintf("%s/%s", q.host, MessagesPath)
}

func (q query) body() (io.Reader, error) {
	var data map[string]interface{} = make(map[string]interface{})
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

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to encode query %w", err)
	}
	return bytes.NewReader(dataJSON), nil
}

func (q query) interval(from, to time.Time) string {
	params := url.Values{}

	params.Add("q", q.query)
	params.Add("fields_in_order", strings.Join(q.fields, ", "))

	params.Add("timerange", "absolute")
	params.Add("from", from.Format(GrayLogDateFormat))
	params.Add("to", to.Format(GrayLogDateFormat))

	return fmt.Sprintf("%s/streams/%s/search?%s", q.host, q.streamid, params.Encode())
}
