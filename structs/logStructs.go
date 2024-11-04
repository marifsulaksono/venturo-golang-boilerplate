package structs

type LogEntry struct {
	URL      string      `json:"url"`
	Path     string      `json:"path"`
	IP       string      `json:"ip"`
	User     interface{} `json:"user"`
	Body     interface{} `json:"body"`
	Response interface{} `json:"response"`
}
