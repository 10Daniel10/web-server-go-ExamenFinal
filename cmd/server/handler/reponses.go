package handler

type ErrorResponse struct {
	Timestamp string   `json:"timestamp"`
	Status    int      `json:"status"`
	Message   string   `json:"message"`
	Path      string   `json:"path"`
	Errors    []string `json:"errors,omitempty"`
} // @name ErrorResponse
