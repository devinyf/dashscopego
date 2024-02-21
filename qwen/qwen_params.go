package qwen

import (
	"log"
)

const (
	DashScopeBaseURL = "https://dashscope.aliyuncs.com/api"
	QwenSubURL       = "/v1/services/aigc/text-generation/generation"
	QwenVLSubURL     = "/v1/services/aigc/multimodal-generation/generation"
)

type ModelQwen string

const (
	QwenTurbo          ModelQwen = "qwen-turbo"
	QwenPlus           ModelQwen = "qwen-plus"
	QwenMax            ModelQwen = "qwen-max"
	QwenMax1201        ModelQwen = "qwen-max-1201"
	QwenMaxLongContext ModelQwen = "qwen-max-longcontext"
)

type Model struct{}

// text-generation only.
func URLQwen() string {
	return DashScopeBaseURL + QwenSubURL
}

// multimodal.
func URLQwenVL() string {
	return DashScopeBaseURL + QwenVLSubURL
}

func ChoseModelQwen(model string) ModelQwen {
	m := Model{}
	switch model {
	case "qwen-turbo":
		return m.QwenTurbo()
	case "qwen-plus":
		return m.QwenPlus()
	case "qwen-max":
		return m.QwenMax()
	case "qwen-max-1201":
		return m.QwenMax1201()
	case "qwen-max-longcontext":
		return m.QwenMaxLongContext()
	default:
		log.Println("target model not found, use default model: qwen-turbo")
		return m.QwenTurbo()
	}
}

func (m *Model) QwenTurbo() ModelQwen {
	return QwenTurbo
}

func (m *Model) QwenPlus() ModelQwen {
	return QwenPlus
}

func (m *Model) QwenMax() ModelQwen {
	return QwenMax
}

func (m *Model) QwenMax1201() ModelQwen {
	return QwenMax1201
}

func (m *Model) QwenMaxLongContext() ModelQwen {
	return QwenMaxLongContext
}
