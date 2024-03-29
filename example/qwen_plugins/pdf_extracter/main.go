package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"

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

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	// 比亚迪财报
	pdfFile := filepath.Join(usr.HomeDir, "Downloads", "financial_results_byd.pdf")
	filePath := fmt.Sprintf("file://%s", pdfFile)

	content := qwen.FileContentList{
		{
			Text: "分析该公司的财报信息, 判断有哪些风险和机会, 并预测后续的业绩发展趋势, 给出投资建议",
		},
		{
			// File: "https://qianwen-res.oss-cn-beijing.aliyuncs.com/QWEN_TECHNICAL_REPORT.pdf",
			File: filePath,
		},
	}

	input := dashscopego.FileInput{
		Messages: []dashscopego.FileMessage{
			{Role: qwen.RoleUser, Content: &content},
		},
	}

	// callback function:  print stream result
	streamCallbackFn := func(ctx context.Context, chunk []byte) error {
		log.Print(string(chunk))
		return nil
	}
	req := &dashscopego.FileRequest{
		Input:       input,
		StreamingFn: streamCallbackFn,
		Plugin:      qwen.Plugins{"pdf_extracter": {}},
	}

	ctx := context.TODO()
	resp, err := cli.CreateFileCompletion(ctx, req)
	if err != nil {
		panic(err)
	}

	log.Println("\nnon-stream result: ")
	log.Println(resp.Output.Choices[0].Message.Content.ToString())
}
