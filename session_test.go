package gograylog

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func Test_SessionJSON(t *testing.T) {
	var vuUnmarshaled ValidUntil
	err := json.Unmarshal([]byte(`"2023-08-28T15:30:00.123-0700"`), &vuUnmarshaled)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	cases := []struct {
		Name    string
		Session Session
	}{
		{
			Name: "Marshal JSON",
			Session: Session{
				Id:         "1",
				Username:   "One",
				ValidUntil: vuUnmarshaled,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.Name, func(t *testing.T) {
			b, err := json.Marshal(tt.Session)
			if err != nil {
				t.Error(err)
			}

			var session Session
			err = json.Unmarshal(b, &session)
			if err != nil {
				t.Error(err)
			}

			if time.Time(tt.Session.ValidUntil).String() != time.Time(session.ValidUntil).String() {
				t.Errorf("time mismatch %s != %s", time.Time(tt.Session.ValidUntil).String(), time.Time(session.ValidUntil).String())
			}
		})
	}

}
