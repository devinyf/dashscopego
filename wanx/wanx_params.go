package wanx

import "fmt"

const (
	ImageSynthesisURI = "/v1/services/aigc/text2image/image-synthesis"
	TaskURI           = "/v1/tasks/%s"
)

type ModelWanx = string

const (
	WanxV1             ModelWanx = "wanx-v1"
	WanxStyleRepaintV1 ModelWanx = "wanx-style-repaint-v1"
	WanxBgGenV2        ModelWanx = "wanx-background-generation-v2"
)

func ImageSynthesisURL(baseURL string) string {
	return baseURL + ImageSynthesisURI
}

func TaskURL(baseURL, taskID string) string {
	return baseURL + fmt.Sprintf(TaskURI, taskID)
}
