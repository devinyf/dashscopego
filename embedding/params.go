package embedding

const (
	embeddingURL = "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
)

type ModelEmbedding = string

const (
	EmbeddingV1      = "text-embedding-v1"
	EmbeddingAsyncV1 = "text-embedding-async-v1"
	EmbeddingV2      = "text-embedding-v2"
	EmbeddingAsyncV2 = "text-embedding-async-v2"
)

type EmbeddingTextType = string

const (
	TypeQuery    EmbeddingTextType = "query"
	TypeDocument EmbeddingTextType = "document"
)
