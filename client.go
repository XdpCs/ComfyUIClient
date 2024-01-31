package comfyUIclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID         string
	httpURL    string
	queueCount int
	webSocket  *WebSocketConnection
	ch         chan *WSMessage
	httpClient *http.Client
}

func NewDefaultClient(serverAddress string, port string) *Client {
	return NewClient(serverAddress, port, &http.Client{Timeout: 10 * time.Second})

}

func NewClient(serverAddress string, port string, httpClient *http.Client) *Client {
	c := &Client{
		ID:         uuid.New().String(),
		httpURL:    "https://" + serverAddress + ":" + port,
		httpClient: httpClient,
		ch:         make(chan *WSMessage),
	}
	c.webSocket = NewDefaultWebSocketConnection("wss://"+serverAddress+":"+port+"/ws?clientId="+c.ID, c)
	return c
}

func (c *Client) IsInitialized() bool {
	return c.webSocket.GetIsConnected() == true
}

func (c *Client) ConnectAndListen() {
	go c.webSocket.ConnectAndListen()
}

func (c *Client) SendTaskStatus(w *WSMessage) error {
	if c.ch != nil {
		c.ch <- w
		return nil
	}
	return errors.New("client not initialized, ch is nil")
}

func (c *Client) GetTaskStatus() chan *WSMessage {
	return c.ch
}

func (c *Client) GetQueueCount() int {
	return c.queueCount
}

func (c *Client) Handle(msg string) error {
	message := &WSMessage{}
	if err := json.Unmarshal([]byte(msg), message); err != nil {
		return fmt.Errorf("json.Unmarshal: error: %w", err)
	}

	switch message.Type {
	case Status:
		s := message.Data.(*WSMessageDataStatus)
		c.queueCount = s.Status.ExecInfo.QueueRemaining
	case ExecutionStart, ExecutionCached, Executing,
		Progress, Executed, ExecutionInterrupted, ExecutionError:
		if err := c.SendTaskStatus(message); err != nil {
			return fmt.Errorf("SendTaskStatus: error: %w", err)
		}
	default:
		return fmt.Errorf("unknown message type: %s, message: %v", message.Type, message)
	}
	return nil
}

// QueuePrompt queues a prompt and starts execution
func (c *Client) QueuePrompt(workflow string) (*QueuePromptResp, error) {
	if !c.IsInitialized() {
		return nil, errors.New("client not initialized")
	}

	if workflow == "" {
		return nil, errors.New("workflow is empty")
	}

	temp := struct {
		ClientID string          `json:"client_id"`
		Prompt   json.RawMessage `json:"prompt"`
	}{
		Prompt:   []byte(workflow),
		ClientID: c.ID,
	}

	reqBody, err := json.Marshal(temp)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: error: %w", err)
	}

	resp, err := c.httpClient.Post(c.httpURL+string(PromptRouter), "application/json", io.NopCloser(bytes.NewReader(reqBody)))
	if err != nil {
		return nil, fmt.Errorf("httpClient.Post: error: %w", err)
	}
	defer resp.Body.Close()
	q := &QueuePromptResp{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}

	if err := json.Unmarshal(body, &q); err != nil {
		return nil, fmt.Errorf("json.NewDecoder: error: %w, resp.Body: %v", err, string(body))
	}

	return q, nil
}

// GetQueueRemaining returns queue remaining
func (c *Client) GetQueueRemaining() (uint64, error) {
	resp, err := c.getJsonUsesRouter(PromptRouter, nil)
	if err != nil {
		return 0, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var queueExecInfo QueueExecInfo
	if err := json.Unmarshal(body, &queueExecInfo); err != nil {
		return 0, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return queueExecInfo.ExecInfo.QueueRemaining, nil
}

// GetEmbeddings returns embeddings
func (c *Client) GetEmbeddings() ([]string, error) {
	resp, err := c.getJsonUsesRouter(EmbeddingsRouter, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var embeddings []string
	if err := json.Unmarshal(body, &embeddings); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return embeddings, nil
}

// GetExtensions returns extensions for frontend
func (c *Client) GetExtensions() ([]string, error) {
	resp, err := c.getJsonUsesRouter(ExtensionsRouter, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var extensions []string
	if err := json.Unmarshal(body, &extensions); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return extensions, nil
}

// GetAllHistories returns all histories
func (c *Client) GetAllHistories() ([]*PromptHistoryItem, error) {
	resp, err := c.getJsonUsesRouter(HistoryRouter, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	return getHistorySlices(resp)
}

// GetHistoryByPromptID returns history info by promptID
func (c *Client) GetHistoryByPromptID(promptID string) (*PromptHistoryItem, error) {
	resp, err := c.getJson(string(HistoryRouter)+"/"+promptID, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	history, err := getHistorySlices(resp)
	if err != nil {
		return nil, fmt.Errorf("getHistorySlices: error: %w", err)
	}
	if len(history) == 0 {
		return nil, nil
	}
	return history[0], nil
}

func getHistorySlices(resp *http.Response) ([]*PromptHistoryItem, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var historyMap map[string]*PromptHistoryMember
	if err := json.Unmarshal(body, &historyMap); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	histories := make([]*PromptHistoryItem, 0, len(historyMap))
	for k, v := range historyMap {
		histories = append(histories, &PromptHistoryItem{
			PromptID:            k,
			PromptHistoryMember: *v,
		})
	}
	return histories, nil
}

// DeleteAllHistories deletes all histories
func (c *Client) DeleteAllHistories() error {
	data := map[string]string{"clear": "clear"}
	_, err := c.postJSONUsesRouter(HistoryRouter, data)
	if err != nil {
		return fmt.Errorf("http.Post: error: %w", err)
	}
	return nil
}

// DeleteHistoryByPromptID deletes history by promptID
func (c *Client) DeleteHistoryByPromptID(promptID string) error {
	data := map[string][]string{"delete": {promptID}}
	_, err := c.postJSONUsesRouter(HistoryRouter, data)
	if err != nil {
		return fmt.Errorf("http.Post: error: %w", err)
	}

	return nil
}

// GetImage returns image byte data
func (c *Client) GetImage(image *DataOutputImages) (*[]byte, error) {
	params := url.Values{}
	params.Add("filename", image.Filename)
	params.Add("subfolder", image.SubFolder)
	params.Add("type", image.Type)
	resp, err := c.getJsonUsesRouter(ViewRouter, params)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	return &body, nil
}

// GetViewMetadata returns view metadata
func (c *Client) GetViewMetadata(folderName string, fileName string) ([]byte, error) {
	if folderName == "" {
		return nil, errors.New("folderName is empty")
	}

	if folderName[0] != '/' {
		folderName = "/" + folderName
	}

	resp, err := c.getJson(string(ViewMetadataRouter)+folderName, url.Values{"filename": {fileName}})
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	return body, nil
}

// GetSystemStats returns system stats
func (c *Client) GetSystemStats() (*SystemStats, error) {
	resp, err := c.getJsonUsesRouter(SystemStatsRouter, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var stats SystemStats
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return &stats, nil
}

// InterruptExecution interrupts execution
func (c *Client) InterruptExecution() error {
	_, err := c.postJSONUsesRouter(InterruptRouter, nil)
	if err != nil {
		return fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	return nil
}

// DeleteAllQueues deletes all prompts in queue
// Delete all prompts in queue with this client sent, or it will not work
func (c *Client) DeleteAllQueues() error {
	data := map[string]string{"clear": "clear"}
	_, err := c.postJSONUsesRouter(QueueRouter, data)
	if err != nil {
		return fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	return nil
}

// DeleteQueueByPromptID deletes prompt in queue by promptID
// You must input promptID with this client sent, or it will not work
func (c *Client) DeleteQueueByPromptID(promptID string) error {
	data := map[string]string{"delete": promptID}
	_, err := c.postJSONUsesRouter(QueueRouter, data)
	if err != nil {
		return fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	return nil
}

func (c *Client) requestJson(method string, endpoint string, values url.Values, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}
	rawURL := c.httpURL + endpoint
	if len(values) != 0 {
		rawURL += values.Encode()
	}
	req, err := http.NewRequest(method, rawURL, io.NopCloser(bytes.NewReader(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.httpClient.Do: %w", err)
	}
	return resp, nil
}

func (c *Client) postJSONUsesRouter(endpoint Router, data interface{}) (*http.Response, error) {
	return c.postJson(string(endpoint), data)
}

func (c *Client) postJson(endpoint string, data interface{}) (*http.Response, error) {
	return c.requestJson(http.MethodPost, endpoint, nil, data)
}

func (c *Client) getJsonUsesRouter(endpoint Router, values url.Values) (*http.Response, error) {
	return c.getJson(string(endpoint), values)
}

func (c *Client) getJson(endpoint string, values url.Values) (*http.Response, error) {
	return c.requestJson(http.MethodGet, endpoint, values, nil)
}
