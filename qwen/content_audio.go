package qwen

type AudioContent struct {
	Audio string `json:"audio,omitempty"`
	Text  string `json:"text,omitempty"`
}

type AudioContentList []AudioContent

func NewAudioContentList() *AudioContentList {
	ac := AudioContentList(make([]AudioContent, 0))
	return &ac
}

func (acList *AudioContentList) ToBytes() []byte {
	if acList == nil || len(*acList) == 0 {
		return []byte("")
	}
	return []byte((*acList)[0].Text)
}

func (acList *AudioContentList) ToString() string {
	if acList == nil || len(*acList) == 0 {
		return ""
	}
	return (*acList)[0].Text
}

func (acList *AudioContentList) SetText(s string) {
	if acList == nil {
		panic("AudioContentList is nil")
	}
	*acList = append(*acList, AudioContent{Text: s})
}

func (acList *AudioContentList) SetAudio(url string) {
	if acList == nil {
		panic("AudioContentList is nil or empty")
	}
	*acList = append(*acList, AudioContent{Audio: url})
}

func (acList *AudioContentList) PopAudioContent() (AudioContent, bool) {
	if acList == nil {
		panic("AudioContentList is nil or empty")
	}

	hasAudio := false
	for i, v := range *acList {
		if v.Audio != "" {
			hasAudio = true
			preSlice := (*acList)[:i]
			if i == len(*acList)-1 {
				*acList = preSlice
			} else {
				postSlice := (*acList)[i+1:]
				*acList = append(*acList, preSlice...)
				*acList = append(*acList, postSlice...)
			}

			return v, hasAudio
		}
	}
	return AudioContent{}, hasAudio
}

func (acList *AudioContentList) AppendText(s string) {
	if acList == nil || len(*acList) == 0 {
		panic("AudioContentList is nil or empty")
	}
	(*acList)[0].Text += s
}

func (acList *AudioContentList) SetImage(_ string) {
	panic("AudioContentList does not support SetImage")
}
