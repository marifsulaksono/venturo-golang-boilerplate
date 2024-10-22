package structs

type Mail struct {
	TargetName  string `json:"target_name"`
	TargetEmail string `json:"target_email"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
}
