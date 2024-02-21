package dashscopego

import (
	"context"
	"strings"

	embedding "github.com/devinyf/dashscopego/embedding"
	httpclient "github.com/devinyf/dashscopego/httpclient"
	qwen "github.com/devinyf/dashscopego/qwen"
	wanx "github.com/devinyf/dashscopego/wanx"
)

type TongyiClient struct {
	Model   string
	token   string
	httpCli httpclient.IHttpClient
}

func NewTongyiClient(model string, token string) *TongyiClient {
	httpcli := httpclient.NewHTTPClient()
	return newTongyiCLientWithHTTPCli(model, token, httpcli)
}

func newTongyiCLientWithHTTPCli(model string, token string, httpcli httpclient.IHttpClient) *TongyiClient {
	return &TongyiClient{
		Model:   model,
		httpCli: httpcli,
		token:   token,
	}
}

// duplicate: CreateCompletion and CreateVLCompletion are the same but with different payload types.
// maybe this can be change in the future.
//
// nolint:lll
func (q *TongyiClient) CreateCompletion(ctx context.Context, payload *qwen.Request[*qwen.TextContent], url string) (*TextQwenResponse, error) {
	payload = paylosdPreCheck(q, payload)
	return genericCompletion(ctx, payload, q.httpCli, url, q.token)
}

//nolint:lll
func (q *TongyiClient) CreateVLCompletion(ctx context.Context, payload *qwen.Request[*qwen.VLContentList], url string) (*VLQwenResponse, error) {
	payload = paylosdPreCheck(q, payload)
	// Uploading URL...
	// TODO: will upload the same image multiple times when this in history messages
	// fmt.Println("upload images...")
	for _, vMsg := range payload.Input.Messages {
		if vMsg.Role == "user" {
			if tmpImageContent, ok := vMsg.Content.PopImageContent(); ok {
				var ossURL string
				var err error
				filepath := tmpImageContent.Image
				// fmt.Println(">>> filepath: ", filepath)
				switch {
				case strings.HasPrefix(filepath, "file://"):
					// fmt.Println(">>> 111111: ", filepath)
					filepath = strings.TrimPrefix(filepath, "file://")
					ossURL, err = qwen.UploadLocalImg(ctx, filepath, payload.Model, q.token)
				case strings.HasPrefix(filepath, "https://") || strings.HasPrefix(filepath, "http://"):
					// fmt.Println(">>> 2222222: ", filepath)
					ossURL, err = qwen.UploadImgFromURL(ctx, filepath, payload.Model, q.token)
				default:
					return nil, ErrImageFilePrefix
				}

				if err != nil {
					return nil, err
				}
				payload.HasUploadOss = true
				// replace the image content with oss url
				// fmt.Printf("after upload, ossURL: %s\n", ossURL)
				vMsg.Content.SetImage(ossURL)
			}
		}
	}

	// =====================
	// fmt.Printf("after upload:  %+v\n", payload.Input.Messages)
	// msgJson, _ := json.Marshal(payload.Input.Messages)
	// fmt.Printf("after upload:  %+v\n", string(msgJson))
	// =====================
	return genericCompletion(ctx, payload, q.httpCli, url, q.token)
}

//nolint:lll
func genericCompletion[T qwen.IQwenContent](ctx context.Context, payload *qwen.Request[T], httpcli httpclient.IHttpClient, url, token string) (*qwen.OutputResponse[T], error) {
	if payload.Model == "" {
		return nil, ErrModelNotSet
	}

	// use streaming if streaming func is set
	if payload.StreamingFn != nil {
		payload.Parameters.SetIncrementalOutput(true)
		return qwen.SendMessageStream(ctx, payload, httpcli, url, token)
	}

	return qwen.SendMessage(ctx, payload, httpcli, url, token)
}

// TODO: intergrate wanx.Request into qwen.IQwenContent(or should rename to ITongyiContent)
//
//nolint:lll
func (q *TongyiClient) CreateImageGeneration(ctx context.Context, payload *wanx.ImageSynthesisRequest) ([]*wanx.ImgBlob, error) {
	if payload.Model == "" {
		if q.Model == "" {
			return nil, ErrModelNotSet
		}
		payload.Model = q.Model
	}
	return wanx.CreateImageGeneration(ctx, payload, q.httpCli, q.token)
}

func (q *TongyiClient) CreateEmbedding(ctx context.Context, r *embedding.Request) ([][]float32, error) {
	resp, err := embedding.CreateEmbedding(ctx, r, q.httpCli, q.token)
	if err != nil {
		return nil, err
	}
	if len(resp.Output.Embeddings) == 0 {
		return nil, ErrEmptyResponse
	}

	embeddings := make([][]float32, 0)
	for i := 0; i < len(resp.Output.Embeddings); i++ {
		embeddings = append(embeddings, resp.Output.Embeddings[i].Embedding)
	}
	return embeddings, nil
}

func paylosdPreCheck[T qwen.IQwenContent](q *TongyiClient, payload *qwen.Request[T]) *qwen.Request[T] {
	if payload.Model == "" {
		payload.Model = q.Model
	}

	if payload.Parameters == nil {
		payload.Parameters = qwen.DefaultParameters()
	}

	return payload
}
