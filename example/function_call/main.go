package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/devinyf/dashscopego"
	"github.com/devinyf/dashscopego/qwen"
)

// 定义工具 tools.
// nolint:all
var tools = []qwen.Tool{
	{
		Type: "function",
		Function: qwen.ToolFunction{
			Name:        "get_current_weather",
			Description: "当你想查询指定城市的天气时非常有用。",
			Parameters: qwen.ToolCallParameter{
				Type: "object",
				Properties: map[string]qwen.PropertieDefinition{
					"location": {
						Type:        "string",
						Description: "城市名称",
					},
				},
			},
			Required: []string{"location"},
		},
	},
}

func getCurrentWeather(cityName string) string {
	return fmt.Sprintf("%v今天有钞票雨。 ", cityName)
}

func main() {
	model := qwen.QwenTurbo
	token := os.Getenv("DASHSCOPE_API_KEY")

	if token == "" {
		panic("token is empty")
	}

	cli := dashscopego.NewTongyiClient(model, token)

	content := qwen.TextContent{Text: "青岛今天的天气怎么样?"}
	messages := []dashscopego.TextMessage{
		{Role: qwen.RoleUser, Content: &content},
	}

NEXT_ROUND:
	input := dashscopego.TextInput{
		Messages: messages,
	}

	req := &dashscopego.TextRequest{
		Input: input,
		Tools: tools,
	}

	ctx := context.TODO()
	resp, err := cli.CreateCompletion(ctx, req)
	if err != nil {
		panic(err)
	}

	log.Println("\nnon-stream result: ")

	if resp.HasToolCallInput() {
		// 需要调用工具
		toolCalls := *resp.Output.Choices[0].Message.ToolCalls

		for _, toolCall := range toolCalls {
			fnName := toolCall.Function.Name
			if fnName == "get_current_weather" {
				argMap := toolCall.Function.GetArguments()
				cityName := argMap["location"]
				toolAnswer := getCurrentWeather(cityName)
				// fmt.Println("tool answer: ", tool_answer)

				toolMessage := dashscopego.TextMessage{
					Role: qwen.RoleTool,
					Content: &qwen.TextContent{
						Text: toolAnswer,
					},
					Name: &fnName,
				}

				// 添加 assistant 的回答
				assistantOutput := resp.Output.Choices[0].Message
				messages = append(messages, assistantOutput)
				// 添加 tool 的回答
				messages = append(messages, toolMessage)

				// 继续下一轮对话
				goto NEXT_ROUND
			}
		}
	}
	// Final result
	// nolint:all
	fmt.Println(resp.Output.Choices[0].Message.Content.ToString())
}
