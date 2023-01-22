package gograylog

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// date format required by graylogs
const GrayLogDateFormat string = "2006-01-02T15:04:05.000Z"

type QueryInterface interface {
	JSON() ([]byte, error)
	Url(string, time.Time, time.Time) string
}

type Query struct {
	QueryString, StreamID string
	Fields                []string
	Limit, Frequency      int
}

// JSON Method converts Query types to appropriate key/value mapping
// then JSON encodes and returns the byte array
func (q Query) JSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["streams"] = []string{q.StreamID}

	data["timerange"] = map[string]string{
		"type":  "relative",
		"range": strconv.Itoa(q.Frequency * 60),
	}

	data["query_string"] = map[string]string{
		"type":         "elasticsearch",
		"query_string": q.QueryString,
	}

	if len(q.Fields) > 0 {
		data["fields_in_order"] = q.Fields
	}

	data["limit"] = q.Limit

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("unable to encode query %w", err)
	}
	return dataJSON, nil
}

// Url method takes a from and to time.Time to combine with Query struct fields to the appropriate
// key value format then converts those to URL param format. The values are url encoded and
// applied to search api endpoint
func (q Query) Url(host string, from, to time.Time) string {
	params := url.Values{}

	params.Add("q", q.QueryString)
	params.Add("fields_in_order", strings.Join(q.Fields, ", "))

	params.Add("timerange", "absolute")
	params.Add("from", from.Format(GrayLogDateFormat))
	params.Add("to", to.Format(GrayLogDateFormat))

	return fmt.Sprintf("%s/streams/%s/search?%s", host, q.StreamID, params.Encode())
}
