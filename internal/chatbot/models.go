package chatbot

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Assistant string `json:"assistant"`
	Reply     string `json:"reply"`
}
