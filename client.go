package comfyUIclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID         string
	baseURL    string
	queueCount int
	webSocket  *WebSocketConnection
	ch         chan *WSMessage
	httpClient *http.Client
}

type EndPoint struct {
	Protocol string
	Address  string
	Port     string
}

func NewEndPoint(protocol, address, port string) *EndPoint {
	return &EndPoint{
		Protocol: protocol,
		Address:  address,
		Port:     port,
	}
}

func (e *EndPoint) String() string {
	if e.Port == "" {
		return e.Protocol + "://" + e.Address
	}
	return e.Protocol + "://" + e.Address + ":" + e.Port
}

func NewDefaultClient(endPoint *EndPoint) *Client {
	return NewClient(endPoint, &http.Client{Timeout: 10 * time.Second})
}

func NewDefaultClientStr(baseURL string) *Client {
	baseURLParsed, _ := url.Parse(baseURL)
	endPoint := NewEndPoint(baseURLParsed.Scheme, baseURLParsed.Hostname(), baseURLParsed.Port())
	return NewClient(endPoint, &http.Client{Timeout: 10 * time.Second})
}

func NewClient(endPoint *EndPoint, httpClient *http.Client) *Client {
	c := &Client{
		ID:         uuid.New().String(),
		baseURL:    endPoint.String(),
		httpClient: httpClient,
		ch:         make(chan *WSMessage),
	}

	if strings.HasPrefix(c.baseURL, "https") {
		endPoint.Protocol = "wss"
	} else {
		endPoint.Protocol = "ws"
	}
	c.webSocket = NewDefaultWebSocketConnection(endPoint.String()+"/ws?clientId="+c.ID, c)
	return c
}

func (c *Client) IsInitialized() bool {
	return c.webSocket.GetIsConnected()
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

// QueuePromptByString queues a prompt and starts execution by workflow which type is string
// workflow must be a json string
// extraDataString must be a json string
func (c *Client) QueuePromptByString(workflow string, extraDataString string) (*QueuePromptResp, error) {
	if !c.IsInitialized() {
		return nil, errors.New("client not initialized")
	}

	if workflow == "" {
		return nil, errors.New("workflow is empty")
	}

	if extraDataString == "" {
		extraDataString = "{}"
	}

	temp := struct {
		ClientID  string          `json:"client_id"`
		Prompt    json.RawMessage `json:"prompt"`
		ExtraData *extraData      `json:"extra_data"`
	}{
		Prompt:   []byte(workflow),
		ClientID: c.ID,
		ExtraData: &extraData{
			ExtraPngInfo: []byte(extraDataString),
		},
	}

	resp, err := c.postJSONUsesRouter(PromptRouter, temp, nil)
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

// QueuePromptByNodes queues a prompt and starts execution by workflow which type is map[string]PromptNode
// extraData must be a json string
func (c *Client) QueuePromptByNodes(nodes map[string]PromptNode, extraDataString string) (*QueuePromptResp, error) {
	if len(nodes) == 0 {
		return nil, errors.New("nodes is empty")
	}

	temp := struct {
		ClientID  string                `json:"client_id"`
		Prompt    map[string]PromptNode `json:"prompt"`
		ExtraData *extraData            `json:"extra_data"`
	}{
		Prompt:   nodes,
		ClientID: c.ID,
		ExtraData: &extraData{
			ExtraPngInfo: []byte(extraDataString),
		},
	}
	return c.queuePrompt(temp)
}

func (c *Client) queuePrompt(temp interface{}) (*QueuePromptResp, error) {
	resp, err := c.postJSONUsesRouter(PromptRouter, temp, nil)
	if err != nil {
		return nil, fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()

	q := &QueuePromptResp{}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}

	if err := json.Unmarshal(body, &q); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}

	return q, nil
}

// GetQueueRemaining returns queue remaining
func (c *Client) GetQueueRemaining() (uint64, error) {
	resp, err := c.getJsonUsesRouter(PromptRouter, nil, nil)
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
	resp, err := c.getJsonUsesRouter(EmbeddingsRouter, nil, nil)
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
	resp, err := c.getJsonUsesRouter(ExtensionsRouter, nil, nil)
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
	resp, err := c.getJsonUsesRouter(HistoryRouter, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	return getHistorySlices(resp)
}

// GetHistoryByPromptID returns history info by promptID
func (c *Client) GetHistoryByPromptID(promptID string) (*PromptHistoryItem, error) {
	resp, err := c.getJson(string(HistoryRouter)+"/"+promptID, nil, nil)
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
	_, err := c.postJSONUsesRouter(HistoryRouter, data, nil)
	if err != nil {
		return fmt.Errorf("http.Post: error: %w", err)
	}
	return nil
}

// DeleteHistoryByPromptID deletes history by promptID
func (c *Client) DeleteHistoryByPromptID(promptID string) error {
	data := map[string][]string{"delete": {promptID}}
	_, err := c.postJSONUsesRouter(HistoryRouter, data, nil)
	if err != nil {
		return fmt.Errorf("http.Post: error: %w", err)
	}

	return nil
}

// GetFile returns file byte data
func (c *Client) GetFile(image *DataOutputFile) (*[]byte, error) {
	params := url.Values{}
	params.Add("filename", image.Filename)
	params.Add("subfolder", image.SubFolder)
	params.Add("type", image.Type)
	resp, err := c.getJsonUsesRouter(ViewRouter, params, nil)
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

	resp, err := c.getJson(string(ViewMetadataRouter)+folderName, url.Values{"filename": {fileName}}, nil)
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
	resp, err := c.getJsonUsesRouter(SystemStatsRouter, nil, nil)
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
	_, err := c.postJSONUsesRouter(InterruptRouter, nil, nil)
	if err != nil {
		return fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	return nil
}

// DeleteAllQueues deletes all prompts in queue
// Delete all prompts in queue with this client sent, or it will not work
func (c *Client) DeleteAllQueues() error {
	data := map[string]string{"clear": "clear"}
	_, err := c.postJSONUsesRouter(QueueRouter, data, nil)
	if err != nil {
		return fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	return nil
}

// DeleteQueueByPromptID deletes prompt in queue by promptID
// You must input promptID with this client sent, or it will not work
func (c *Client) DeleteQueueByPromptID(promptID string) error {
	data := map[string]string{"delete": promptID}
	_, err := c.postJSONUsesRouter(QueueRouter, data, nil)
	if err != nil {
		return fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	return nil
}

// GetObjectInfos returns node infos in workflow
func (c *Client) GetObjectInfos() (map[string]*NodeObject, error) {
	resp, err := c.getJsonUsesRouter(ObjectInfoRouter, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var objectInfos map[string]*NodeObject
	if err := json.Unmarshal(body, &objectInfos); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return objectInfos, nil
}

// GetObjectInfoByNodeName returns node info by nodeName
func (c *Client) GetObjectInfoByNodeName(name string) (*NodeObject, error) {
	resp, err := c.getJson(string(ObjectInfoRouter)+"/"+name, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJson: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var objectInfos map[string]*NodeObject
	if err := json.Unmarshal(body, &objectInfos); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return objectInfos[name], nil
}

// GetQueueInfo returns queue info
func (c *Client) GetQueueInfo() (*QueueInfo, error) {
	resp, err := c.getJsonUsesRouter(QueueRouter, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("c.getJsonUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}
	var queueInfo *QueueInfo
	if err := json.Unmarshal(body, &queueInfo); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}
	return queueInfo, nil
}

func (c *Client) uploadFile(router Router, reader io.Reader, fileName string, overwrite bool, filetype ImageType, subFolder string) (*UploadFile, error) {
	requestBody, headers, err := createUploadRequest(reader, fileName, overwrite, filetype, subFolder)
	if err != nil {
		return nil, fmt.Errorf("createUploadRequest: error: %w", err)
	}

	resp, err := c.postMultiPartUsesRouter(router, requestBody, headers)
	if err != nil {
		return nil, fmt.Errorf("c.postJSONUsesRouter: error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: error: %w", err)
	}

	var u *UploadFile
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, fmt.Errorf("json.Unmarshal: error: %w, resp.Body: %v", err, string(body))
	}

	return u, nil
}

// UploadImage uploads image
func (c *Client) UploadImage(reader io.Reader, fileName string, overwrite bool, filetype ImageType, subFolder string) (*UploadFile, error) {
	return c.uploadFile(UploadImageRouter, reader, fileName, overwrite, filetype, subFolder)
}

// UploadMask uploads mask image
func (c *Client) UploadMask(reader io.Reader, fileName string, overwrite bool, filetype ImageType, subFolder string) (*UploadFile, error) {
	return c.uploadFile(UploadMaskRouter, reader, fileName, overwrite, filetype, subFolder)
}

func createUploadRequest(reader io.Reader, fileName string, overwrite bool, filetype ImageType, subFolder string) (*bytes.Buffer, map[string]string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)
	defer writer.Close()
	formFile, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return nil, nil, err
	}
	_, err = io.Copy(formFile, reader)
	if err != nil {
		return nil, nil, fmt.Errorf("io.Copy: %w", err)
	}

	if err := writer.WriteField("overwrite", fmt.Sprintf("%v", overwrite)); err != nil {
		return nil, nil, fmt.Errorf("writer.WriteField: overwrite %v overwrite %w", overwrite, err)
	}

	if err := writer.WriteField("type", fmt.Sprintf("%v", filetype)); err != nil {
		return nil, nil, fmt.Errorf("writer.WriteField: type %v error: %w", filetype, err)
	}

	if subFolder != "" {
		if err := writer.WriteField("subfolder", fmt.Sprintf("%v", subFolder)); err != nil {
			return nil, nil, fmt.Errorf("writer.WriteField: subfolder %v error: %w", subFolder, err)
		}
	}

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	return &requestBody, headers, nil
}

func (c *Client) makeRequest(method, router string, values url.Values, data interface{}, headers map[string]string, contentType string) (*http.Response, error) {
	var req *http.Request
	var err error

	rawURL := c.baseURL + router
	if len(values) != 0 {
		rawURL += "?" + values.Encode()
	}

	if data != nil {
		switch contentType {
		case "application/json":
			jsonData, err := json.Marshal(data)
			if err != nil {
				return nil, fmt.Errorf("json.Marshal: %w", err)
			}
			req, err = http.NewRequest(method, rawURL, io.NopCloser(bytes.NewReader(jsonData)))
			if err != nil {
				return nil, fmt.Errorf("http.NewRequest: %w", err)
			}
		case "multipart/form-data":
			buf := data.(*bytes.Buffer)
			req, err = http.NewRequest(method, rawURL, io.NopCloser(buf))
			if err != nil {
				return nil, fmt.Errorf("http.NewRequest: %w", err)
			}
		default:
			return nil, fmt.Errorf("unsupported content type: %s", contentType)
		}
	} else {
		req, err = http.NewRequest(method, rawURL, nil)
		if err != nil {
			return nil, fmt.Errorf("http.NewRequest: %w", err)
		}
	}

	// Don't change the order
	req.Header.Set("Content-Type", contentType)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("c.httpClient.Do: %w", err)
	}
	return resp, nil
}

func (c *Client) requestJson(method, router string, values url.Values, data interface{}, headers map[string]string) (*http.Response, error) {
	return c.makeRequest(method, router, values, data, headers, "application/json")
}

func (c *Client) requestMultiPart(method, router string, values url.Values, data *bytes.Buffer, headers map[string]string) (*http.Response, error) {
	return c.makeRequest(method, router, values, data, headers, "multipart/form-data")
}

func (c *Client) postMultiPartUsesRouter(router Router, data *bytes.Buffer, headers map[string]string) (*http.Response, error) {
	return c.requestMultiPart(http.MethodPost, string(router), nil, data, headers)
}

func (c *Client) postJSONUsesRouter(router Router, data interface{}, headers map[string]string) (*http.Response, error) {
	return c.postJson(string(router), data, headers)
}

func (c *Client) postJson(router string, data interface{}, headers map[string]string) (*http.Response, error) {
	return c.requestJson(http.MethodPost, router, nil, data, headers)
}

func (c *Client) getJsonUsesRouter(router Router, values url.Values, headers map[string]string) (*http.Response, error) {
	return c.getJson(string(router), values, headers)
}

func (c *Client) getJson(router string, values url.Values, headers map[string]string) (*http.Response, error) {
	return c.requestJson(http.MethodGet, router, values, nil, headers)
}
