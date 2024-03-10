package entities

type Request struct {
	Source      string            `json:"src"`
	Destination string            `json:"dst"`
	Protocol    string            `json:"proto"`
	Method      string            `json:"method"`
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
}
