### dashscopego

阿里云平台 dashscope api 的 golang 封装 (非官方)

[开通DashScope并创建API-KEY](https://help.aliyun.com/zh/dashscope/developer-reference/activate-dashscope-and-create-an-api-key)

Install:
```sh
go get -u github.com/devinyf/dashscopego
```

#### Examples:
* [通义千问](#通义千问)
* [通义千问VL(视觉理解模型)](#通义千问VL视觉理解模型)
* [通义千问Audio(音频语言模型)](#通义千问Audio音频语言模型)
* [通义万相(图像生成)](#通义万相图像生成)
* [Paraformer(语音识别)](#Paraformer语音识别)
* 模型插件调用 TODO
* langchaingo Agent TODO

开发中...

### 通义千问
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

### 通义万相(图像生成)
- [x] 文本生成图像
- [ ] 人像风格重绘
- [ ] 图像背景生成
```go
func main() {
	model := wanx.WanxV1
	token := os.Getenv("DASHSCOPE_API_KEY")
	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	req := &wanx.ImageSynthesisRequest{
		// Model: "wanx-v1",
		Model: model,
		Input: wanx.ImageSynthesisInput{
			Prompt: "画一只松鼠",
		},
		Params: wanx.ImageSynthesisParams{
			N: 1,
		},
		Download: true // 从 URL 下载图片
	}
	ctx := context.TODO()

	imgBlobs, err := cli.CreateImageGeneration(ctx, req)
	if err != nil {
		panic(err)
	}

	for _, blob := range imgBlobs {
		// blob.Data 会在 request 中设置了 Download: true 时下载
		// 否则使用 blob.ImgURL
		saveImg2Desktop(blob.ImgType, blob.Data)
	}
}

func saveImg2Desktop(fileType string, data []byte) {
	buf := bytes.NewBuffer(data)
	img, _, err := image.Decode(buf)
	if err != nil {
		log.Fatal(err)
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	f, err := os.Create(filepath.Join(usr.HomeDir, "Desktop", "wanx_image.png"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
```

### 通义千问VL(视觉理解模型)
 * Image 也可以直接使用 图片本地路径 或 图片URL链接的, 参照了 dashscope python 库的实现步骤 临时上传到 oss
 * 其中上传图片到 oss 的步骤 在开发文档中还没有看到HTTP调用的例子, 所以后续可能会做变更
```go
func main() {
	model := qwen.QwenVLPlus
	token := os.Getenv("DASHSCOPE_API_KEY")

	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	sysContent := qwen.VLContentList{
		{
			Text: "You are a helpful assistant.",
		},
	}
	userContent := qwen.VLContentList{
		{
			Text: "用唐诗体描述一下这张图片中的内容",
		},
		{
            // 使用 图片URL链接
            // Image: "https://pic.ntimg.cn/20140113/8800276_184351657000_2.jpg",
            // 本地图片
            // Image: "file:///Users/xxxx/xxxx.png",
            // 官方文档的例子, oss 下载
			Image: "https://dashscope.oss-cn-beijing.aliyuncs.com/images/dog_and_girl.jpeg",
		},
	}

	input := dashscopego.VLInput{
		Messages: []dashscopego.VLMessage{
			{Role: "system", Content: &sysContent},
			{Role: "user", Content: &userContent},
		},
	}

	// (可选 SSE开启)需要流式输出时 通过该 Callback Function 获取结果
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}
	req := &dashscopego.VLRequest{
		Input:       input,
		StreamingFn: streamCallbackFn,
	}

	ctx := context.TODO()
	resp, err := cli.CreateVLCompletion(ctx, req)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nnon-stream result: ")
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString())
}
```

### 通义千问Audio(音频语言模型)
* 同 QwenVL, 如果使用本地音频文件会临时上传 oss, 之后可能会有变动
```go
func main() {
	model := qwen.QwenAudioTurbo
	token := os.Getenv("DASHSCOPE_API_KEY")

	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	sysContent := qwen.AudioContentList{
		{
			Text: "You are a helpful assistant.",
		},
	}
	userContent := qwen.AudioContentList{
		{
			Text: "该段对话表达了什么观点? 详细分析该讲话者的语气,展现出什么样的情绪", //nolint:gosmopolitan
		},
		{
			// 使用本地音频文件
			// Audio: "file:///Users/xxx/Desktop/hello_world_female2.wav",
			// 官方文档中的例子
			Audio: "https://dashscope.oss-cn-beijing.aliyuncs.com/audios/2channel_16K.wav",
		},
	}

	input := dashscopego.AudioInput{
		Messages: []dashscopego.AudioMessage{
			{Role: "system", Content: &sysContent},
			{Role: "user", Content: &userContent},
		},
	}

	// callback function:  print stream result
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		log.Print(string(chunk))
		return nil
	}
	req := &dashscopego.AudioRequest{
		Input:       input,
		StreamingFn: streamCallbackFn,
	}

	ctx := context.TODO()
	resp, err := cli.CreateAudioCompletion(ctx, req)
	if err != nil {
		panic(err)
	}

	log.Println("\nnon-stream result: ")
	log.Println(resp.Output.Choices[0].Message.Content.ToString())
}
```

### Paraformer(语音识别)
- [x] 实时语音识别API
- [ ] 录音文件识别API

Experimental:
* 开发文档中 还没有看到 HTTP调用说明, 参照 dashscope python 库中的步骤实现, 将来可能会有变更
* 参数中的: SampleRate 好像目前仅支持 16000, 使用真实录音要留意录音设备的 sample_rate 是与之否匹配
```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/paraformer"
)

func main() {
	model := paraformer.ParaformerRealTimeV1
	token := os.Getenv("DASHSCOPE_API_KEY")
	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}

	headerPara := paraformer.ReqHeader{
		Streaming: "duplex",
		TaskID:    paraformer.GenerateTaskID(),
		Action:    "run-task",
	}

	payload := paraformer.PayloadIn{
		Parameters: paraformer.Parameters{
			// seems like only support 16000 sample-rate.
			SampleRate: 16000,
			Format:     "pcm",
		},
		Input:     map[string]interface{}{},
		Task:      "asr",
		TaskGroup: "audio",
		Function:  "recognition",
	}

	req := &paraformer.Request{
		Header:      headerPara,
		Payload:     payload,
		StreamingFn: streamCallbackFn,
	}

	// 声音获取 实际使用时请替换成实时音频流.
	voiceReader := readAudioFromDesktop()

	reader := bufio.NewReader(voiceReader)

	cli.CreateSpeechToTextGeneration(context.TODO(), req, reader)

	// 等待语音识别结果输出
	time.Sleep(5 * time.Second)
}

// 读取音频文件中的录音 模拟实时语音流. 这里下载的官方文档中的示例音频文件.
// `https://dashscope.oss-cn-beijing.aliyuncs.com/samples/audio/paraformer/hello_world_male2.wav`.
func readAudioFromDesktop() *bufio.Reader {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	voiceFilePath := filepath.Join(usr.HomeDir, "Desktop", "hello_world_female2.wav")
	f, err := os.OpenFile(voiceFilePath, os.O_RDONLY, 0640)
	if err != nil {
		panic(err)
	}
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(f)
	return reader
}
```