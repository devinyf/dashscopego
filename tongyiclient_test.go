package dashscopego

import (
	"context"
	"os"
	"strings"
	"testing"

	httpclient "github.com/devinyf/dashscopego/httpclient"
	qwen "github.com/devinyf/dashscopego/qwen"
	wanx "github.com/devinyf/dashscopego/wanx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTongyiClient(t *testing.T, model string) *TongyiClient {
	t.Helper()
	token := os.Getenv("DASHSCOPE_API_KEY")

	cli := NewTongyiClient(model, token)
	if cli.token == "" {
		t.Skip("token is empty")
	}
	return cli
}

func newMockClient(t *testing.T, model string, ctrl *gomock.Controller, f mockFn) *TongyiClient {
	t.Helper()

	mockHTTPCli := httpclient.NewMockIHttpClient(ctrl)
	fackToken := ""

	f(mockHTTPCli)

	qwenCli := newTongyiCLientWithHTTPCli(model, fackToken, mockHTTPCli)
	return qwenCli
}

type mockFn func(mockHTTPCli *httpclient.MockIHttpClient)

func TestBasic(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	cli := newTongyiClient(t, "")

	text := qwen.TextContent{Text: "Hello"}
	input := TextInput{
		Messages: []TextMessage{
			{Role: "user", Content: &text},
		},
	}

	req := &TextRequest{
		Model: "qwen-turbo",
		Input: input,
	}

	resp, err := cli.CreateCompletion(ctx, req, qwen.URLQwen())

	require.NoError(t, err)
	assert.Regexp(t, "hello|hi|how|assist", resp.Output.Choices[0].Message.Content.ToString())
}

func TestStreamingChunk(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	cli := newTongyiClient(t, "qwen-turbo")

	content := qwen.TextContent{Text: "Hello"}

	input := TextInput{
		Messages: []TextMessage{
			{Role: "user", Content: &content},
		},
	}

	output := ""
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		output += string(chunk)
		return nil
	}

	req := &TextRequest{
		// Model: "qwen-turbo",
		Input:       input,
		StreamingFn: streamCallbackFn,
	}
	resp, err := cli.CreateCompletion(ctx, req, qwen.URLQwen())

	require.NoError(t, err)
	assert.Regexp(t, "hello|hi|how|assist", resp.Output.Choices[0].Message.Content.ToString())
	assert.Regexp(t, "hello|hi|how|assist", output)
}

func TestVLBasic(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	cli := newTongyiClient(t, "")

	sysContent := qwen.VLContentList{
		{
			Text: "You are a helpful assistant.",
		},
	}
	userContent := qwen.VLContentList{
		{
			Text: "describe the image",
		},
		{
			Image: "https://dashscope.oss-cn-beijing.aliyuncs.com/images/dog_and_girl.jpeg",
		},
	}

	input := VLInput{
		Messages: []VLMessage{
			{Role: "system", Content: &sysContent},
			{Role: "user", Content: &userContent},
		},
	}

	req := &VLRequest{
		Model: "qwen-vl-plus",
		Input: input,
	}

	req.Parameters = qwen.DefaultParameters()

	resp, err := cli.CreateVLCompletion(ctx, req, qwen.URLQwenVL())

	require.NoError(t, err)
	assert.Regexp(t, "dog|person|individual|woman|girl", resp.Output.Choices[0].Message.Content.ToString())
}

func TestVLStreamChund(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	cli := newTongyiClient(t, "")

	sysContent := qwen.VLContentList{
		{
			Text: "You are a helpful assistant.",
		},
	}
	userContent := qwen.VLContentList{
		{
			Text: "describe the image",
		},
		{
			Image: "https://dashscope.oss-cn-beijing.aliyuncs.com/images/dog_and_girl.jpeg",
			// Image: "https://pic.ntimg.cn/20140113/8800276_184351657000_2.jpg",
		},
	}

	input := VLInput{
		Messages: []VLMessage{
			{Role: "system", Content: &sysContent},
			{Role: "user", Content: &userContent},
		},
	}

	output := ""
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		output += string(chunk)
		return nil
	}

	req := &VLRequest{
		Model:       "qwen-vl-plus",
		Input:       input,
		StreamingFn: streamCallbackFn,
	}

	resp, err := cli.CreateVLCompletion(ctx, req, qwen.URLQwenVL())

	require.NoError(t, err)
	assert.Equal(t, output, resp.Output.Choices[0].Message.Content.ToString())
	assert.Regexp(t, "dog|person|individual|woman|girl", strings.ToLower(output))
}

func TestImageGeneration(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	cli := newTongyiClient(t, "wanx-v1")

	req := &wanx.ImageSynthesisRequest{
		Model: "wanx-v1",
		Input: wanx.ImageSynthesisInput{
			Prompt: "A beautiful sunset",
		},
		Download: true,
	}

	imgBlobs, err := cli.CreateImageGeneration(ctx, req)
	require.NoError(t, err)
	require.NotEmpty(t, imgBlobs)

	for _, blob := range imgBlobs {
		assert.NotEmpty(t, blob.Data)
		assert.NotEmpty(t, blob.ImgURL)
		assert.Equal(t, "image/png", blob.ImgType)
	}
}

func TestMockStreamingChunk(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := newMockClient(t, "qwen-turbo", ctrl, _mockAsyncFunc)

	output := ""
	text := qwen.TextContent{Text: "Hello"}
	input := TextInput{
		Messages: []TextMessage{
			{Role: "user", Content: &text},
		},
	}

	req := &TextRequest{
		Input: input,
		StreamingFn: func(ctx context.Context, chunk []byte) error {
			output += string(chunk)
			return nil
		},
	}

	mockURL := ""
	resp, err := cli.CreateCompletion(ctx, req, mockURL)

	require.NoError(t, err)

	assert.Equal(t, "Hello! How can I assist you today?", resp.Output.Choices[0].Message.Content.ToString())
	assert.Equal(t, "Hello! How can I assist you today?", output)
}

func TestMockBasic(t *testing.T) {
	t.Parallel()
	ctx := context.TODO()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cli := newMockClient(t, "qwen-turbo", ctrl, _mockSyncFunc)
	text := qwen.TextContent{Text: "Hello"}
	input := TextInput{
		Messages: []TextMessage{
			{Role: "user", Content: &text},
		},
	}

	req := &TextRequest{
		Input: input,
	}

	mockURL := ""
	resp, err := cli.CreateCompletion(ctx, req, mockURL)

	require.NoError(t, err)

	assert.Equal(t, "Hello! This is a mock message.", resp.Output.Choices[0].Message.Content.ToString())
	assert.Equal(t, "mock-ac55-9fd3-8326-8415cbdf5683", resp.RequestID)
	assert.Equal(t, 15, resp.Usage.TotalTokens)
}

func _mockAsyncFunc(mockHTTPCli *httpclient.MockIHttpClient) {
	MockStreamData := []string{
		`id:1`,
		`event:result`,
		`:HTTP_STATUS/200`,
		`data:{
			"output": {
				"choices": [{
					"message": {
						"content": "Hello! How",
						"role": "assistant"
					},
					"finish_reason": "null"
				}]
			},
			"usage": {
				"total_tokens": 9,
				"input_tokens": 6,
				"output_tokens": 3
			},
			"request_id": "95bea986-ac55-9fd3-8326-8415cbdf5683"
		}`,
		`    `,
		`id:2`,
		`event:result`,
		`:HTTP_STATUS/200`,
		`data:{
			"output": {
				"choices": [{
					"message": {
						"content": " can I assist you today?",
						"role": "assistant"
					},
					"finish_reason": "null"
				}]
			},
			"usage": {
				"total_tokens": 15,
				"input_tokens": 6,
				"output_tokens": 9
			},
			"request_id": "95bea986-ac55-9fd3-8326-8415cbdf5683"
		}`,
		`    `,
		`id:3`,
		`event:result`,
		`:HTTP_STATUS/200`,
		`data:{
			"output": {
				"choices": [{
					"message": {
						"content": "",
						"role": "assistant"
					},
					"finish_reason": "stop"
				}]
			},
			"usage": {
				"total_tokens": 15,
				"input_tokens": 6,
				"output_tokens": 9
			},
			"request_id": "95bea986-ac55-9fd3-8326-8415cbdf5683"
		}`,
	}

	ctx := context.TODO()

	_rawStreamOutChannel := make(chan string, 500)

	mockHTTPCli.EXPECT().
		PostSSE(ctx, gomock.Any(), gomock.Any(), gomock.Any()).Return(_rawStreamOutChannel, nil)
	go func() {
		for _, line := range MockStreamData {
			_rawStreamOutChannel <- line
		}
		close(_rawStreamOutChannel)
	}()
}

func _mockSyncFunc(mockHTTPCli *httpclient.MockIHttpClient) {
	ctx := context.TODO()

	text := qwen.TextContent{Text: "Hello! This is a mock message."}

	mockResp := TextQwenResponse{
		Output: qwen.Output[*qwen.TextContent]{
			Choices: []qwen.Choice[*qwen.TextContent]{
				{
					Message: TextMessage{
						Content: &text,
						Role:    "assistant",
					},
					FinishReason: "stop",
				},
			},
		},
		Usage: struct {
			TotalTokens  int `json:"total_tokens"`
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		}{
			TotalTokens:  15,
			InputTokens:  6,
			OutputTokens: 9,
		},
		RequestID: "mock-ac55-9fd3-8326-8415cbdf5683",
	}
	mockHTTPCli.EXPECT().
		Post(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		SetArg(3, mockResp).
		Return(nil)
}
