package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/tools/dashscope/wanx"

	"github.com/tmc/langchaingo/llms/tongyi"
	"github.com/tmc/langchaingo/tools"
)

/*
	使用通义千问 调用 通义万象 实现文生图功能
	这个例子展示了如何使用 langchaingo Agent 来调用通义千问(llm) 和 通义万象(tool)

	*实验性功能 还没有并入主分支*
	*注意在 go.mod 中 replace 了 langchiango 到 fork 功能分支*
*/
func main() {
	if err := agentExample(); err != nil {
		panic(err)
	}
}

func agentExample() error {
	modelOpt := tongyi.WithModel("qwen-turbo")
	keyOpt := tongyi.WithToken(os.Getenv("DASHSCOPE_API_KEY"))

	llm, err := tongyi.New(modelOpt, keyOpt)
	if err != nil {
		return err
	}

	wanxDescOpt := wanx.WithDescription(wanx.WanxDescriptionCN) // 切换中文描述 默认是英文
	wanxTool := wanx.NewTongyiWanx(wanxDescOpt)

	agentTools := []tools.Tool{
		tools.Calculator{},
		wanxTool,
	}

	callbackHandler := callbacks.NewFinalStreamHandler()
	// 打印全部输出信息, 包括中间 agent 调用工具的过程
	callbackHandler.PrintOutput = true

	agent := agents.NewOneShotAgent(llm, agentTools, agents.WithCallbacksHandler(callbackHandler))
	executor := agents.NewExecutor(agent)

	// 设置 streaming 结果的回调函数
	fn := func(ctx context.Context, chunk []byte) {
		fmt.Printf("%s", string(chunk))
	}

	// 获取 streaming 输出
	callbackHandler.ReadFromEgress(context.Background(), fn)

	fmt.Printf("\n")

	question := "画一个武松打虎, 帮我生成提示词, 要写实的画风, 把老虎画的萌一点"
	answer, err := chains.Run(context.Background(), executor, question)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n<Final Answer>: ", answer)
	_ = answer

	return err
}
