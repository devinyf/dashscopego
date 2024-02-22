package paraformer

import "context"

const ParaformerWSURL = "wss://dashscope.aliyuncs.com/api-ws/v1/inference"

type Parameters struct {
	SampleRate int    `json:"sample_rate"`
	Format     string `json:"format"`
}

type StreamingFunc func(ctx context.Context, chunk []byte) error

type Request struct {
	Header      ReqHeader     `json:"header"`
	Payload     PayloadIn     `json:"payload"`
	StreamingFn StreamingFunc `json:"-"`
}

type ReqHeader struct {
	Streaming string `json:"streaming"`
	TaskID    string `json:"task_id"`
	Action    string `json:"action"`
}

type PayloadIn struct {
	Model      string                 `json:"model"`
	Parameters Parameters             `json:"parameters"`
	Input      map[string]interface{} `json:"input"`
	Task       string                 `json:"task"`
	TaskGroup  string                 `json:"task_group"`
	Function   string                 `json:"function"`
}

// ---------
// type Word struct {
// 	BeginTime   int    `json:"begin_time"`
// 	EndTime     int    `json:"end_time"`
// 	Text        string `json:"text"`
// 	Punctuation string `json:"punctuation"`
// }

type Sentence struct {
	BeginTime int    `json:"begin_time"`
	EndTime   int    `json:"end_time"`
	Text      string `json:"text"` // full text
	// Words     []Word `json:"words"`
}

type Output struct {
	Sentence Sentence `json:"sentence"`
}

type Usage struct {
	Duration int `json:"duration"`
}

type PayloadOut struct {
	Output Output `json:"output"`
	Usage  Usage  `json:"usage"`
}

type Attributes struct{}

type Header struct {
	TaskID     string     `json:"task_id"`
	Event      string     `json:"event"`
	Attributes Attributes `json:"attributes"`
}

type RecognitionResult struct {
	Header  Header     `json:"header"`
	Payload PayloadOut `json:"payload"`
}
