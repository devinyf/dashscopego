package dashscopego

import (
	"github.com/devinyf/dashscopego/qwen"
)

type (
	TextInput  = qwen.Input[*qwen.TextContent]
	VLInput    = qwen.Input[*qwen.VLContentList]
	AudioInput = qwen.Input[*qwen.AudioContentList]
	FileInput  = qwen.Input[*qwen.FileContentList]

	TextRequest  = qwen.Request[*qwen.TextContent]
	VLRequest    = qwen.Request[*qwen.VLContentList]
	AudioRequest = qwen.Request[*qwen.AudioContentList]
	FileRequest  = qwen.Request[*qwen.FileContentList]

	TextQwenResponse  = qwen.OutputResponse[*qwen.TextContent]
	VLQwenResponse    = qwen.OutputResponse[*qwen.VLContentList]
	AudioQwenResponse = qwen.OutputResponse[*qwen.AudioContentList]
	FileQwenResponse  = qwen.OutputResponse[*qwen.FileContentList]

	TextMessage  = qwen.Message[*qwen.TextContent]
	VLMessage    = qwen.Message[*qwen.VLContentList]
	AudioMessage = qwen.Message[*qwen.AudioContentList]
	FileMessage  = qwen.Message[*qwen.FileContentList]
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
