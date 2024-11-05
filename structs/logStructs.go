package structs

import "time"

type LogEntry struct {
	Datetime time.Time   `json:"datetime"`
	URL      string      `json:"url"`
	Method   string      `json:"method"`
	IP       string      `json:"ip"`
	User     interface{} `json:"user"`
	Body     interface{} `json:"body"`
	Response interface{} `json:"response"`
}
