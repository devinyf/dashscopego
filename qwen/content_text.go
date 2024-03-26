package qwen

import (
	"encoding/json"
)

// TextConent is used for text-generation only.
type TextContent struct {
	Text  string
	IsRaw bool // 临时，如果是raw数据，直接返回Text
}

var _ IQwenContentMethods = &TextContent{}

func NewTextContent() *TextContent {
	t := TextContent{
		Text: "",
	}
	return &t
}

func (t *TextContent) ToBytes() []byte {
	return []byte(t.Text)
}

func (t *TextContent) ToString() string {
	return t.Text
}

func (t *TextContent) SetText(text string) {
	if t == nil {
		panic("TextContent is nil")
	}
	t.Text = text
}

func (t *TextContent) AppendText(text string) {
	if t == nil {
		panic("TextContent is nil")
	}
	t.Text += text
}

func (t *TextContent) SetImage(_ string) {
	panic("TextContent does not support SetImage")
}

func (t *TextContent) SetAudio(_ string) {
	panic("TextContent does not support SetAudio")
}

// redifine MarshalJSON and UnmarshalJSON.
func (t TextContent) MarshalJSON() ([]byte, error) {
	if t.IsRaw {
		return []byte(t.Text), nil
	}

	return json.Marshal(t.Text)
}

func (t *TextContent) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &t.Text)
}
