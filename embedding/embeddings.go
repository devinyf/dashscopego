package embedding

import (
	"context"

	"github.com/devinyf/dashscopego/httpclient"
)

//nolint:lll
func CreateEmbedding(ctx context.Context, req *Request, cli httpclient.IHttpClient, token string) (*Response, error) {
	if req.Model == "" {
		req.Model = EmbeddingV2
	}
	if req.Params.TextType == "" {
		req.Params.TextType = TypeDocument
	}

	resp := Response{}
	tokenOption := httpclient.WithTokenHeaderOption(token)
	err := cli.Post(ctx, embeddingURL, req, &resp, tokenOption)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
