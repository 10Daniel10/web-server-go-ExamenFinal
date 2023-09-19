package handler

type ErrorResponse struct {
	Timestamp string   `json:"timestamp"`
	Status    int      `json:"status"`
	Message   string   `json:"message"`
	Path      string   `json:"path"`
	Errors    []string `json:"errors,omitempty"`
} // @name Error

type PageResponse struct {
	Page     int         `json:"page"`
	Size     int         `json:"size"`
	Total    int         `json:"total"`
	Content  interface{} `json:"content"`
	NextPage string      `json:"next_page"`
	PrevPage string      `json:"prev_page"`
}
