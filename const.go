package comfyUIclient

type WsMessageType string

const (
	Status               WsMessageType = "status"
	Progress             WsMessageType = "progress"
	Executed             WsMessageType = "executed"
	Executing            WsMessageType = "executing"
	ExecutionStart       WsMessageType = "execution_start"
	ExecutionError       WsMessageType = "execution_error"
	ExecutionCached      WsMessageType = "execution_cached"
	ExecutionInterrupted WsMessageType = "execution_interrupted"
)

type Router string

const (
	PromptRouter       Router = "/prompt"
	HistoryRouter      Router = "/history"
	ViewRouter         Router = "/view"
	ViewMetadataRouter Router = "/view_metadata"
	EmbeddingsRouter   Router = "/embeddings"
	ExtensionsRouter   Router = "/extensions"
	SystemStatsRouter  Router = "/system_stats"
	InterruptRouter    Router = "/interrupt"
	QueueRouter        Router = "/queue"
	ObjectInfoRouter   Router = "/object_info"
	UploadImageRouter  Router = "/upload/image"
	UploadMaskRouter   Router = "/upload/mask"
)

type TaskStatusType = WsMessageType

type ImageType string

const (
	InputImageType  ImageType = "input"
	TempImageType   ImageType = "temp"
	OutputImageType ImageType = "output"
)
