package qwen

// qwen(text-generation) and qwen-vl(multi-modal) have different data format
// so define generic interfaces for them.
type IQwenContent interface {
	*TextContent | *VLContentList | *AudioContentList
	IQwenContentMethods
}

type IQwenContentMethods interface {
	ToBytes() []byte
	ToString() string
	SetText(text string)
	AppendText(text string)
}

// TextConent is used for text-generation only.
type TextContent string

func NewTextContent() *TextContent {
	t := TextContent("")
	return &t
}

func (t *TextContent) ToBytes() []byte {
	str := *t
	return []byte(str)
}

func (t *TextContent) ToString() string {
	str := *t
	return string(str)
}

func (t *TextContent) SetText(text string) {
	*t = TextContent(text)
}

func (t *TextContent) AppendText(text string) {
	str := *t
	*t = TextContent(string(str) + text)
}
