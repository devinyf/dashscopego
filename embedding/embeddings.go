package embedding

import (
	"context"

	"github.com/devinyf/dashscopego/httpclient"
)

const (
	embeddingURL          = "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
	defaultEmbeddingModel = "text-embedding-v1"
)

type Request struct {
	Model string `json:"model"`
	Input struct {
		Texts []string `json:"texts"`
	} `json:"input"`
	Params struct {
		TextType string `json:"text_type"` // query or document
	} `json:"parameters"`
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

type Response struct {
	Output Output `json:"output"`
}

//nolint:lll
func CreateEmbedding(ctx context.Context, req *Request, cli httpclient.IHttpClient, token string) (*Response, error) {
	if req.Model == "" {
		req.Model = defaultEmbeddingModel
	}
	if req.Params.TextType == "" {
		req.Params.TextType = "document"
	}

	resp := Response{}
	tokenOption := httpclient.WithTokenHeaderOption(token)
	err := cli.Post(ctx, embeddingURL, req, &resp, tokenOption)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
