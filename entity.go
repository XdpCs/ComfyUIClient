package comfyUIclient

// @Title        entity.go
// @Description
// @Create       XdpCs 2024-01-31 00:06
// @Update       XdpCs 2024-01-31 00:06

import "encoding/json"

// SystemStats contains a system info and gpu infos
type SystemStats struct {
	System  *System `json:"system"`
	Devices []*GPU  `json:"devices"`
}

// System contains system info
type System struct {
	OS             string `json:"os"`
	PythonVersion  string `json:"python_version"`
	EmbeddedPython bool   `json:"embedded_python"`
}

// GPU contains gpu info
type GPU struct {
	Name           string `json:"name"`
	Type           string `json:"type"`
	Index          int    `json:"index"`
	VRAMTotal      int64  `json:"vram_total"`
	VRAMFree       int64  `json:"vram_free"`
	TorchVRAMTotal int64  `json:"torch_vram_total"`
	TorchVRAMFree  int64  `json:"torch_vram_free"`
}

// QueueExecInfo contains queue remaining number
type QueueExecInfo struct {
	ExecInfo struct {
		QueueRemaining uint64 `json:"queue_remaining"`
	} `json:"exec_info"`
}

// QueuePromptResp contains prompt id, number and node errors
type QueuePromptResp struct {
	PromptID   string                 `json:"prompt_id"`
	Number     int                    `json:"number"`
	NodeErrors map[string]interface{} `json:"node_errors"`
}

// DataOutputImages export data address, name and type
type DataOutputImages struct {
	Filename  string `json:"filename"`
	SubFolder string `json:"subfolder"`
	Type      string `json:"type"`
}

// PromptHistoryMember is part of prompt history
type PromptHistoryMember struct {
	WorkFlow json.RawMessage                      `json:"prompt"`
	Outputs  map[string]PromptHistoryMemberImages `json:"outputs"`
}

type PromptHistoryMemberImages struct {
	Images *[]DataOutputImages `json:"images"`
}

// PromptHistoryItem contains prompt id, WorkFlow, output info
type PromptHistoryItem struct {
	PromptID string
	PromptHistoryMember
}
