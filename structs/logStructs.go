package structs

type LogEntry struct {
	URL      string      `json:"url"`
	Path     string      `json:"path"`
	IP       string      `json:"ip"`
	User     interface{} `json:"user"`
	Response interface{} `json:"response"`
}
