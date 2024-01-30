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
	PromptRouter      Router = "/prompt"
	HistoryRouter     Router = "/history"
	ViewRouter        Router = "/view"
	EmbeddingsRouter  Router = "/embeddings"
	ExtensionsRouter  Router = "/extensions"
	SystemStatsRouter Router = "/system_stats"
)

type TaskStatusType = WsMessageType
