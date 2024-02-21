package tongyiclient

import (
	"github.com/devinyf/dashscopego/qwen"
)

type (
	TextInput        = qwen.Input[*qwen.TextContent]
	VLInput          = qwen.Input[*qwen.VLContentList]
	TextRequest      = qwen.Request[*qwen.TextContent]
	VLRequest        = qwen.Request[*qwen.VLContentList]
	TextQwenResponse = qwen.OutputResponse[*qwen.TextContent]
	VLQwenResponse   = qwen.OutputResponse[*qwen.VLContentList]
	TextMessage      = qwen.Message[*qwen.TextContent]
	VLMessage        = qwen.Message[*qwen.VLContentList]
)

func NewQwenMessage[T qwen.IQwenContent](role string, content T) *qwen.Message[T] {
	if content == nil {
		panic("content is nil")
	}

	return &qwen.Message[T]{
		Role:    role,
		Content: content,
	}
}
