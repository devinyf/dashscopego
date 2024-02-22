package paraformer

type ModelParaformer = string

const (
	// detect from file.
	ParaformerV1    ModelParaformer = "paraformer-v1"
	Paraformer8KV1  ModelParaformer = "paraformer-8k-v1"
	ParaformerMtlV1 ModelParaformer = "paraformer-mtl-v1"
	// real time voice.
	ParaformerRealTimeV1   ModelParaformer = "paraformer-realtime-v1"
	ParaformerRealTime8KV1 ModelParaformer = "paraformer-realtime-8k-v1"
)
