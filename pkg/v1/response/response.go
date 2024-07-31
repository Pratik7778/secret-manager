package response

type SuccessResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ListResponse struct {
	Secrets    []string `json:"secrets"`
	Query      string   `json:"query"`
	PageNumber int      `json:"page_number"`
	PageSize   int      `json:"page_size"`
	Total      int      `json:"total"`
}

type ViewSecretResponse struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
