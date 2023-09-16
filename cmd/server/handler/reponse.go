package handler

type ResponseError struct {
	Timestamp string   `json:"timestamp"`
	Status    int      `json:"status"`
	Message   string   `json:"message"`
	Path      string   `json:"path"`
	Errors    []string `json:"errors,omitempty"`
}
