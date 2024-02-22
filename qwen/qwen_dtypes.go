package qwen

import (
	"context"
)

type Parameters struct {
	ResultFormat      string  `json:"result_format,omitempty"`
	Seed              int     `json:"seed,omitempty"`
	MaxTokens         int     `json:"max_tokens,omitempty"`
	TopP              float64 `json:"top_p,omitempty"`
	TopK              int     `json:"top_k,omitempty"`
	Temperature       float64 `json:"temperature,omitempty"`
	EnableSearch      bool    `json:"enable_search,omitempty"`
	IncrementalOutput bool    `json:"incremental_output,omitempty"`
}

func NewParameters() *Parameters {
	return &Parameters{}
}

const DefaultTemperature = 1.0

func DefaultParameters() *Parameters {
	q := Parameters{}
	q.
		SetResultFormat("message").
		SetTemperature(DefaultTemperature)

	return &q
}

func (p *Parameters) SetResultFormat(value string) *Parameters {
	p = p.tryInit()
	p.ResultFormat = value
	return p
}

func (p *Parameters) SetSeed(value int) *Parameters {
	p = p.tryInit()
	p.Seed = value
	return p
}

func (p *Parameters) SetMaxTokens(value int) *Parameters {
	p = p.tryInit()
	p.MaxTokens = value
	return p
}

func (p *Parameters) SetTopP(value float64) *Parameters {
	p = p.tryInit()
	p.TopP = value
	return p
}

func (p *Parameters) SetTopK(value int) *Parameters {
	p = p.tryInit()
	p.TopK = value
	return p
}

func (p *Parameters) SetTemperature(value float64) *Parameters {
	p.tryInit()
	p.Temperature = value
	return p
}

func (p *Parameters) SetEnableSearch(value bool) *Parameters {
	p = p.tryInit()
	p.EnableSearch = value
	return p
}

func (p *Parameters) SetIncrementalOutput(value bool) *Parameters {
	p = p.tryInit()
	p.IncrementalOutput = value
	return p
}

func (p *Parameters) tryInit() *Parameters {
	if p == nil {
		p = &Parameters{}
	}
	return p
}

type Message[T IQwenContent] struct {
	Role    string `json:"role"`
	Content T      `json:"content"`
}

type Input[T IQwenContent] struct {
	Messages []Message[T] `json:"messages"`
}

type StreamingFunc func(ctx context.Context, chunk []byte) error

type Request[T IQwenContent] struct {
	Model      string      `json:"model"`
	Input      Input[T]    `json:"input"`
	Parameters *Parameters `json:"parameters,omitempty"`
	// streaming callback function.
	StreamingFn StreamingFunc `json:"-"`
	// qwen-vl model need to upload image to oss for recognition.
	HasUploadOss bool `json:"-"`
}

func (q *Request[T]) SetModel(value string) *Request[T] {
	q.Model = value
	return q
}

func (q *Request[T]) SetInput(value Input[T]) *Request[T] {
	q.Input = value
	return q
}

func (q *Request[T]) SetParameters(value *Parameters) *Request[T] {
	q.Parameters = value
	return q
}

func (q *Request[T]) SetStreamingFunc(fn func(ctx context.Context, chunk []byte) error) *Request[T] {
	q.StreamingFn = fn
	return q
}

type StreamOutput[T IQwenContent] struct {
	ID         string            `json:"id"`
	Event      string            `json:"event"`
	HTTPStatus int               `json:"http_status"`
	Output     OutputResponse[T] `json:"output"`
	Err        error             `json:"error"`
}

type Choice[T IQwenContent] struct {
	Message      Message[T] `json:"message"`
	FinishReason string     `json:"finish_reason"`
}

// new version response format.
type Output[T IQwenContent] struct {
	Choices []Choice[T] `json:"choices"`
}

type Usage struct {
	TotalTokens  int `json:"total_tokens"`
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type OutputResponse[T IQwenContent] struct {
	Output    Output[T] `json:"output"`
	Usage     Usage     `json:"usage"`
	RequestID string    `json:"request_id"`
	// ErrMsg    string `json:"error_msg"`
}

func (t *OutputResponse[T]) GetChoices() []Choice[T] {
	return t.Output.Choices
}

func (t *OutputResponse[T]) GetUsage() Usage {
	return t.Usage
}

func (t *OutputResponse[T]) GetRequestID() string {
	return t.RequestID
}
