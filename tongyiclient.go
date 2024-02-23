package dashscopego

import (
	"bufio"
	"context"
	"log"
	"strings"

	embedding "github.com/devinyf/dashscopego/embedding"
	httpclient "github.com/devinyf/dashscopego/httpclient"
	"github.com/devinyf/dashscopego/paraformer"
	"github.com/devinyf/dashscopego/qwen"
	"github.com/devinyf/dashscopego/wanx"
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

	for _, vMsg := range payload.Input.Messages {
		tmpImageContent, hasImg := vMsg.Content.PopImageContent()
		if hasImg && vMsg.Role == "user" {
			filepath := tmpImageContent.Image

			ossURL, hasUploadOss, err := checkIfNeedUploadFile(ctx, filepath, payload.Model, q.token)
			if err != nil {
				return nil, err
			}
			if hasUploadOss {
				payload.HasUploadOss = true
			}
			vMsg.Content.SetImage(ossURL)
		}
	}

	return genericCompletion(ctx, payload, q.httpCli, url, q.token)
}

//nolint:lll
func (q *TongyiClient) CreateAudioCompletion(ctx context.Context, payload *qwen.Request[*qwen.AudioContentList], url string) (*AudioQwenResponse, error) {
	payload = paylosdPreCheck(q, payload)
	for _, acMsg := range payload.Input.Messages {
		tmpAudioContent, hasAudio := acMsg.Content.PopAudioContent()

		if hasAudio && acMsg.Role == "user" {
			filepath := tmpAudioContent.Audio

			ossURL, hasUploadOss, err := checkIfNeedUploadFile(ctx, filepath, payload.Model, q.token)
			if err != nil {
				return nil, err
			}

			if hasUploadOss {
				payload.HasUploadOss = true
			}
			acMsg.Content.SetAudio(ossURL)
		}
	}

	return genericCompletion(ctx, payload, q.httpCli, url, q.token)
}

func checkIfNeedUploadFile(ctx context.Context, filepath string, model, token string) (string, bool, error) {
	var err error
	var ossURL string
	var hasUploadOss bool
	switch {
	case strings.Contains(filepath, "dashscope.oss"):
		// 使用了官方案例中的格式(https://dashscope.oss...).
		ossURL = filepath
	case strings.HasPrefix(filepath, "oss://"):
		// 已经在 oss 中的不必上传.
		ossURL = filepath
	case strings.HasPrefix(filepath, "file://"):
		// 本地文件.
		filepath = strings.TrimPrefix(filepath, "file://")
		ossURL, err = qwen.UploadLocalFile(ctx, filepath, model, token)
		hasUploadOss = true
	case strings.HasPrefix(filepath, "https://") || strings.HasPrefix(filepath, "http://"):
		// 文件的 URL 链接.
		ossURL, err = qwen.UploadFileFromURL(ctx, filepath, model, token)
		hasUploadOss = true
	}

	return ossURL, hasUploadOss, err
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

/*
func (q *TongyiClient) CreateVoiceFileToTextGeneration(ctx context.Context, request *paraformer.Request) (any, error) {
	if request.Payload.Model == "" {
		if q.Model == "" {
			return nil, ErrModelNotSet
		}
		request.Payload.Model = q.Model
	}

	return
	panic("not implemented")
}
*/

func (q *TongyiClient) CreateSpeechToTextGeneration(ctx context.Context, request *paraformer.Request, reader *bufio.Reader) error {
	if request.Payload.Model == "" {
		if q.Model == "" {
			return ErrModelNotSet
		}
		request.Payload.Model = q.Model
	}

	wsCli, err := paraformer.ConnRecognitionClient(request, q.token)
	if err != nil {
		return err
	}

	// handle response by stream callback
	go paraformer.HandleRecognitionResult(ctx, wsCli, request.StreamingFn)

	for {
		// this buf can not be reused,
		// otherwise the data will be overwritten, voice became disorder.
		buf := make([]byte, 1024)
		n, errRead := reader.Read(buf)
		if n == 0 {
			break
		}
		if errRead != nil {
			log.Printf("read line error: %v\n", errRead)
			err = errRead
			return err
		}

		paraformer.SendRadioData(wsCli, buf)
	}

	return nil
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
