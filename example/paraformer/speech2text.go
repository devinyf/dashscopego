package main

import (
	"context"
	"log"
	"os"

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

	// callback function:  print stream result
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		log.Print(string(chunk))
		return nil
	}

	headerPara := paraformer.ReqHeader{
		Streaming: "duplex",
		TaskID:    paraformer.GenerateTaskID(),
		Action:    "run-task",
	}

	payload := paraformer.PayloadIn{
		Parameters: paraformer.Parameters{
			// seems like only support 16000 sample-rate now
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

	cli.CreateSpeechToTextGeneration(context.TODO(), req)

}
