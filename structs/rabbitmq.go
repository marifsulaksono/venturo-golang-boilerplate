package structs

type RabbitMQDefaultPayload struct {
	Route string      `json:"command"`
	Param interface{} `json:"param"`
	Data  interface{} `json:"data"`
}

type MessagePayload struct {
	Id         int64       `json:"Id"`
	Command    string      `json:"Command"`
	Time       string      `json:"Time"`
	ModuleId   string      `json:"ModuleId"`
	Properties interface{} `json:"Properties"`
	Signature  string      `json:"Signature"`
	Data       interface{} `json:"Data"`
}

type Request struct {
	ID         int64       `json:"Id"`
	Command    string      `json:"Command"`
	Time       string      `json:"Time"`
	ModuleId   string      `json:"ModuleId"`
	Properties interface{} `json:"Properties"`
	Signature  string      `json:"Signature"`
	Data       interface{} `json:"Data"`
}

type ReqProp struct {
	Offset   int    `json:"Skip"`
	Limit    int    `json:"Take"`
	OrderBy  string `json:"OrderBy"`
	OrderSeq string `json:"OrderSeq"`
}

type Filter struct {
	Id     int32
	Limit  int32
	Offset int32
	Filter string
}
