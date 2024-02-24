package qwen

// VLContentList is used for multi-modal generation.
type VLContentList []VLContent

var _ IQwenContentMethods = &VLContentList{}

type VLContent struct {
	Image string `json:"image,omitempty"`
	Text  string `json:"text,omitempty"`
}

var _ IQwenContentMethods = &VLContentList{}

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

	hasImage := false
	for i, v := range *vlist {
		if v.Image != "" {
			hasImage = true
			preSlice := (*vlist)[:i]
			if i == len(*vlist)-1 {
				*vlist = preSlice
			} else {
				postSlice := (*vlist)[i+1:]
				*vlist = append(*vlist, preSlice...)
				*vlist = append(*vlist, postSlice...)
			}

			return v, hasImage
		}
	}
	return VLContent{}, hasImage
}

func (vlist *VLContentList) AppendText(s string) {
	if vlist == nil || len(*vlist) == 0 {
		panic("VLContentList is nil or empty")
	}
	(*vlist)[0].Text += s
}

func (vlist *VLContentList) SetAudio(_ string) {
	panic("VLContentList does not support SetAudio")
}
