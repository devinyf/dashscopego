### dashscopego

阿里云平台 dashscope api 的 golang 封装 (非官方)

#### examples:
* [通义千问](#通义千问)
* [通义千问VL](#通义千问vl视觉理解大模型)
* [通义万相](#通义万相文生图)

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

	content := qwen.TextContent("讲个冷笑话")

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: "user", Content: &content},
		},
	}

	// callback function:  print stream result
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}
	req := &dashscopego.TextRequest{
		Input:       input,
		StreamingFn: streamCallbackFn,
	}

	ctx := context.TODO()
	resp, err := cli.CreateCompletion(ctx, req, qwen.URLQwen())
	if err != nil {
		panic(err)
	}

	fmt.Println("\nnon-stream result: ")
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString())
}
```

### 通义万相(文生图)
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
	}
	ctx := context.TODO()

	imgBlobs, err := cli.CreateImageGeneration(ctx, req)
	if err != nil {
		panic(err)
	}

	for _, blob := range imgBlobs {
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

### 通义千问VL(视觉理解大模型)
 * P.S. 直接使用 图片本地路径 或 图片URL链接 目前还没有看官方 HTTP接口文档, 这里暂时模拟了 dashscope python 库的实现步骤, 后续可能会做变更
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
            // 官方文档的例子, oss 下载
			Image: "https://dashscope.oss-cn-beijing.aliyuncs.com/images/dog_and_girl.jpeg",
            // 使用 图片URL链接
            // Image: "https://pic.ntimg.cn/20140113/8800276_184351657000_2.jpg",
            // 本地图片
            // Image: "file:///Users/xxxx/xxxx.png",
		},
	}

	input := dashscopego.VLInput{
		Messages: []dashscopego.VLMessage{
			{Role: "system", Content: &sysContent},
			{Role: "user", Content: &userContent},
		},
	}

	// callback function:  print stream result
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		fmt.Print(string(chunk))
		return nil
	}
	req := &dashscopego.VLRequest{
		Input:       input,
		StreamingFn: streamCallbackFn,
	}

	ctx := context.TODO()
	resp, err := cli.CreateVLCompletion(ctx, req, qwen.URLQwenVL())
	if err != nil {
		panic(err)
	}

	fmt.Println("\nnon-stream result: ")
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString())
}
```
### TODO...
* 通义千问Audio(音频语言大模型)
* Paraformer (语音识别转文字)
* 模型插件调用(申请了试用,目前还没通过审核...)