package comfyUIclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
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

type QueuePromptResp struct {
	PromptID   string                 `json:"prompt_id"`
	Number     int                    `json:"number"`
	NodeErrors map[string]interface{} `json:"node_errors"`
}

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

func (c *Client) DeleteAllHistory() error {
	data := map[string]string{"clear": "clear"}
	_, err := c.postJSON(HistoryRouter, data)
	if err != nil {
		return fmt.Errorf("http.Post: error: %w", err)
	}
	return nil
}

func (c *Client) DeleteHistory(promptID string) error {
	data := map[string][]string{"delete": {promptID}}
	_, err := c.postJSON(HistoryRouter, data)
	if err != nil {
		return fmt.Errorf("http.Post: error: %w", err)
	}

	return nil
}

func (c *Client) GetImage(image *DataOutputImages) (*[]byte, error) {
	params := url.Values{}
	params.Add("filename", image.Filename)
	params.Add("subfolder", image.SubFolder)
	params.Add("type", image.Type)
	resp, err := c.httpClient.Get(fmt.Sprintf("%s%s?%s", c.httpURL, ViewRouter, params.Encode()))
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	return &body, nil
}

func (c *Client) postJSON(endpoint Router, data interface{}) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	resp, err := c.httpClient.Post(c.httpURL+string(endpoint), "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, fmt.Errorf("http.Post: %w, url: %v", err, c.httpURL+string(endpoint))
	}
	defer resp.Body.Close()

	return resp, nil
}
