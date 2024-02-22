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
	// Uploading URL...
	// TODO: will upload the same image multiple times when this in history messages
	// fmt.Println("upload images...")
	for _, vMsg := range payload.Input.Messages {
		if vMsg.Role == "user" {
			if tmpImageContent, hasImg := vMsg.Content.PopImageContent(); hasImg {
				var ossURL string
				var err error
				filepath := tmpImageContent.Image
				switch {
				case strings.Contains(filepath, "dashscope.oss"):
					ossURL = filepath
				case strings.HasPrefix(filepath, "file://"):
					filepath = strings.TrimPrefix(filepath, "file://")
					ossURL, err = qwen.UploadLocalImg(ctx, filepath, payload.Model, q.token)
				case strings.HasPrefix(filepath, "https://") || strings.HasPrefix(filepath, "http://"):
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

	return genericCompletion(ctx, payload, q.httpCli, url, q.token)
}

// TODO: COPY CreateVLCompletion 稍后整理消除重复代码
//
//nolint:lll
func (q *TongyiClient) CreateAudioCompletion(ctx context.Context, payload *qwen.Request[*qwen.AudioContentList], url string) (*AudioQwenResponse, error) {
	payload = paylosdPreCheck(q, payload)
	// Uploading URL...
	// fmt.Println("upload audio...")
	// /*
	for _, acMsg := range payload.Input.Messages {
		if acMsg.Role == "user" {
			if tmpImageContent, hasAudio := acMsg.Content.PopAudioContent(); hasAudio {
				var ossURL string
				var err error
				filepath := tmpImageContent.Audio

				// TODO: 使用了 Image 的上传函数，需要修改
				switch {
				case strings.Contains(filepath, "dashscope.oss"):
					ossURL = filepath
				case strings.HasPrefix(filepath, "file://"):
					filepath = strings.TrimPrefix(filepath, "file://")
					ossURL, err = qwen.UploadLocalImg(ctx, filepath, payload.Model, q.token)
				case strings.HasPrefix(filepath, "https://") || strings.HasPrefix(filepath, "http://"):
					ossURL, err = qwen.UploadImgFromURL(ctx, filepath, payload.Model, q.token)
				default:
					return nil, ErrImageFilePrefix
				}

				if err != nil {
					return nil, err
				}
				// TODO: QwenAudio 不能携带 X-DashScope-OssResourceResolve = enable, 报错问题待排查
				// payload.HasUploadOss = true
				// replace the image content with oss url
				// fmt.Printf("after upload, ossURL: %s\n", ossURL)
				acMsg.Content.SetAudio(ossURL)
			}
		}
	}
	// */

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
