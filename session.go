package gograylog

import (
	"errors"
	"time"
)

const (
	GraylogLayout string = "2006-01-02T15:04:05.999-0700"
)

var (
	errMissingSession error = errors.New("missing session data")
	errInvalidSession error = errors.New("invalid session")
)

type ValidUntil time.Time

type Session struct {
	Id         string     `json:"session_id"`
	Username   string     `json:"username"`
	ValidUntil ValidUntil `json:"valid_until"`
}

func (s Session) IsValid() error {
	if s.Id == "" {
		return errMissingSession
	}

	if time.Now().After(time.Time(s.ValidUntil)) {
		return errInvalidSession
	}

	return nil
}

func (v *ValidUntil) UnmarshalJSON(d []byte) error {
	if string(d) == "null" {
		return nil
	}
	d = d[1 : len(d)-1]

	parsedTime, err := time.Parse(GraylogLayout, string(d))
	if err != nil {
		return err
	}

	*v = ValidUntil(parsedTime)

	return nil
}

func (v ValidUntil) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(v).Format(GraylogLayout)), nil
}
