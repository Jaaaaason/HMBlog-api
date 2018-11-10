package handler

// ErrRes error data structure for response
type errRes struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
