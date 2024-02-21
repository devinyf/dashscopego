package qwen

// qwen(text-generation) and qwen-vl(multi-modal) have different data format
// so define generic interfaces for them.
type IQwenContent interface {
	*TextContent | *VLContentList
	IQwenContentMethods
}

type IQwenContentMethods interface {
	ToBytes() []byte
	ToString() string
	SetText(text string)
	SetImage(url string)
	PopImageContent() (VLContent, bool)
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

func (t *TextContent) SetImage(_ string) {
	panic("text-generation only model: can not use SetImage for TextContent")
}

func (vlist *TextContent) PopImageContent() (VLContent, bool) {
	panic("text-generation only model: can not use PopImage for TextContent")
}

func (t *TextContent) AppendText(text string) {
	str := *t
	*t = TextContent(string(str) + text)
}

// VLContentList is used for multi-modal generation.
type VLContentList []VLContent

type VLContent struct {
	Image string `json:"image,omitempty"`
	Text  string `json:"text,omitempty"`
}

func NewVLContentList() *VLContentList {
	vl := VLContentList(make([]VLContent, 0))
	return &vl
}

func (vlist *VLContentList) ToBytes() []byte {
	if vlist == nil || len(*vlist) == 0 {
		return []byte("")
	}
	return []byte((*vlist)[0].Text)
}

func (vlist *VLContentList) ToString() string {
	if vlist == nil || len(*vlist) == 0 {
		return ""
	}
	return (*vlist)[0].Text
}

func (vlist *VLContentList) SetText(s string) {
	if vlist == nil {
		panic("VLContentList is nil")
	}
	*vlist = append(*vlist, VLContent{Text: s})
}

func (vlist *VLContentList) SetImage(url string) {
	if vlist == nil {
		panic("VLContentList is nil or empty")
	}
	*vlist = append(*vlist, VLContent{Image: url})
}

func (vlist *VLContentList) PopImageContent() (VLContent, bool) {
	if vlist == nil {
		panic("VLContentList is nil or empty")
	}

	isOk := false
	for i, v := range *vlist {
		if v.Image != "" {
			isOk = true
			preSlice := (*vlist)[:i]
			if i == len(*vlist)-1 {
				*vlist = preSlice
			} else {
				postSlice := (*vlist)[i+1:]
				*vlist = append(preSlice, postSlice...)
			}

			return v, isOk
		}
	}
	return VLContent{}, isOk
}

func (vlist *VLContentList) AppendText(s string) {
	if vlist == nil || len(*vlist) == 0 {
		panic("VLContentList is nil or empty")
	}
	(*vlist)[0].Text += s
}
