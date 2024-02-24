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
	SetImage(url string)
	SetAudio(url string)
}
