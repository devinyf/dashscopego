package main

import (
	"context"
	"log"
	"os"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
)

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
			Text: "这段音频在说什么", //nolint:gosmopolitan
		},
		{
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
	resp, err := cli.CreateAudioCompletion(ctx, req, qwen.URLQwenAudio())
	if err != nil {
		panic(err)
	}

	log.Println("\nnon-stream result: ")
	log.Println(resp.Output.Choices[0].Message.Content.ToString())
}
