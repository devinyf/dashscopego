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

type TextType = string

const (
	TypeQuery    TextType = "query"
	TypeDocument TextType = "document"
)
