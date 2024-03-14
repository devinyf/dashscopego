package embedding

type Request struct {
	Model string `json:"model"`
	Input struct {
		Texts []string `json:"texts"`
	} `json:"input"`
	Params struct {
		TextType string `json:"text_type"` // query or document
	} `json:"parameters"`
}

type Response struct {
	Output Output `json:"output"`
}

type Embedding struct {
	TextIndex int       `json:"text_index"`
	Embedding []float32 `json:"embedding"`
}

type Output struct {
	Embeddings []Embedding `json:"embeddings"`
	Usgae      struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	RequestID string `json:"request_id"`
}
