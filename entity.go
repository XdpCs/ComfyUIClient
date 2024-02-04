package comfyUIclient

import (
	"encoding/json"
	"fmt"
)

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

// DataOutputFile export data address, name and type
type DataOutputFile struct {
	Filename  string `json:"filename"`
	SubFolder string `json:"subfolder"`
	Type      string `json:"type"`
}

// PromptHistoryMember is part of prompt history
type PromptHistoryMember struct {
	NodeInfo *NodeInfo                            `json:"prompt"`
	Outputs  map[string]PromptHistoryMemberImages `json:"outputs"`
}

type PromptHistoryMemberImages struct {
	Images *[]DataOutputFile `json:"images"`
}

// PromptHistoryItem contains prompt id, WorkFlow, output info
type PromptHistoryItem struct {
	PromptID string
	PromptHistoryMember
}

// PromptNode is the data that inputs into ComfyUI
type PromptNode struct {
	Inputs    map[string]interface{} `json:"inputs"`
	ClassType string                 `json:"class_type"`
}

// NodeObject is a part of workflow
type NodeObject struct {
	Input        *NodeObjectInput `json:"input"`
	Output       []string         `json:"output"`
	OutputIsList []bool           `json:"output_is_list"`
	OutputName   []string         `json:"output_name"`
	Name         string           `json:"name"`
	DisplayName  string           `json:"display_name"`
	Description  string           `json:"description"`
	Category     string           `json:"category"`
	OutputNode   bool             `json:"output_node"`
}

// NodeObjectInput exposes the input information of a node
type NodeObjectInput struct {
	Required map[string]interface{} `json:"required"`
	Optional map[string]interface{} `json:"optional,omitempty"`
}

// QueueInfo exposes the queue info
type QueueInfo struct {
	QueueRunning []*NodeInfo `json:"queue_running"`
	QueuePending []*NodeInfo `json:"queue_pending"`
}

// NodeInfo contains the node info
type NodeInfo struct {
	Num           uint64
	PromptID      string
	Prompt        map[string]PromptNode `json:"prompt"`
	ExtraData     json.RawMessage       // extra data is just for user's custom data
	OutputNodeIDs []string
}

// UploadFile export data address, name and type
type UploadFile struct {
	Filename  string `json:"name"`
	SubFolder string `json:"subfolder"`
	Type      string `json:"type"`
}

func (n *NodeInfo) UnmarshalJSON(data []byte) error {
	var temp []json.RawMessage
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Ensure the length of the array is as expected
	if len(temp) != 5 {
		return fmt.Errorf("unexpected JSON array length for NodeInfo")
	}

	// Extract values from the array
	if err := json.Unmarshal(temp[0], &n.Num); err != nil {
		return err
	}

	if err := json.Unmarshal(temp[1], &n.PromptID); err != nil {
		return err
	}

	if err := json.Unmarshal(temp[2], &n.Prompt); err != nil {
		return err
	}

	n.ExtraData = temp[3]

	if err := json.Unmarshal(temp[4], &n.OutputNodeIDs); err != nil {
		return err
	}

	return nil
}

type extraData struct {
	ExtraPngInfo json.RawMessage `json:"extra_pnginfo"`
}
