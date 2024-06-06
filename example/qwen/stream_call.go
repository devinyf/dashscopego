package main

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

	content := qwen.TextContent{Text: "tell me a joke"}

	input := dashscopego.TextInput{
		Messages: []dashscopego.TextMessage{
			{Role: qwen.RoleUser, Content: &content},
		},
	}

	// callback function:  print stream result
	streamCallbackFn := func(_ context.Context, chunk []byte) error {
		fmt.Print(string(chunk)) //nolint:all
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

	fmt.Println("\nnon-stream result: ")                           //nolint:all
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString()) //nolint:all
}
