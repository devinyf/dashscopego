### dashscopego

阿里云平台 dashscope api 的 golang 封装 (非官方)

[开通DashScope并创建API-KEY](https://help.aliyun.com/zh/dashscope/developer-reference/activate-dashscope-and-create-an-api-key)

Install:
```sh
go get -u github.com/devinyf/dashscopego
```


### Examples:
#### 通义千问
- [x] [大语言模型](./example/qwen/stream_call.go)
- [x] [千问VL(视觉理解模型)](./example/qwen_vl/stream_call.go)
- [x] [千问Audio(音频语言模型)](./example/qwen_audio/stream_call.go)
#### 通义万相
- [x] [文本生成图像](./example/wanx/img_generation.go)
- [ ] 人像风格重绘
- [ ] 图像背景生成
#### Paraformer(语音识别)
- [x] [实时语音识别](./example/paraformer/speech2text.go)
- [ ] 录音文件识别
#### 模型插件调用
- [x] [pdf解析](./example/qwen_plugins/pdf_extracter/main.go)
- [ ] 计算器
- [ ] 图片生成
- [x] [Python代码解释器](./example/qwen_plugins/code_interpreter/main.go)
- [ ] 自定义plugin-example
#### langchaingo-Agent 
- TODO...

开发中...


```go
import (
	"context"
	"fmt"
	"os"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
)

func main() {
	model := qwen.QwenTurbo
	token := os.Getenv("DASHSCOPE_API_KEY")

	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	content := qwen.TextContent{Text: "讲个冷笑话"}

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: "user", Content: &content},
		},
	}

	// (可选 SSE开启) 需要流式输出时 通过该 Callback Function 获取实时显示的结果
	// 开启 SSE 时的 request_id/finish_reason/token usage 等信息在调用完成统一返回(resp)
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}
	req := &dashscopego.TextRequest{
		Input:       input,
		StreamingFn: streamCallbackFn,
	}

	ctx := context.TODO()
	resp, err := cli.CreateCompletion(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nnon-stream result: ")
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString())

	// request_id, finish_reason, token usage
	fmt.Println(resp.RequestID)
	fmt.Println(resp.Output.Choices[0].FinishReason)
	fmt.Println(resp.Usage.TotalTokens)
}
```
