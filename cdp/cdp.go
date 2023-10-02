package cdp

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"os/exec"
)

// Page represents a CDP page.
type Page struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Type  string `json:"type"`
}

// Client provides methods to interact with the CDP.
type Client struct {
	conn *websocket.Conn
}

// NewClient creates a new CDP client.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn}
}

// StartChromium starts a headless Chromium instance and returns the WebSocket Debugger URL.
func StartChromium() (string, error) {
	// Command to start Chromium in headless mode
	cmd := exec.Command("chromium", "--headless", "--disable-gpu", "--remote-debugging-port=9222")
	cmd.Stderr = nil // Set to nil to silence stderr

	// Start the Chromium process
	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to start Chromium: %v", err)
	}

	// Get the WebSocket Debugger URL
	devToolsURL, err := getWebSocketDebuggerURL()
	if err != nil {
		return "", err
	}

	return devToolsURL, nil
}

// getWebSocketDebuggerURL makes a request to obtain the WebSocket Debugger URL
func getWebSocketDebuggerURL() (string, error) {
	resp, err := http.Get("http://127.0.0.1:9222/json/version")
	if err != nil {
		return "", fmt.Errorf("failed to retrieve WebSocket Debugger URL: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	fmt.Println(result.WebSocketDebuggerURL)

	return result.WebSocketDebuggerURL, nil
}

// CreatePage creates a new page and returns its ID.
func (c *Client) CreatePage(url string) (*Page, error) {
	// Send a command to create a new target (page).
	createTargetCmd := map[string]interface{}{
		"id":     1,
		"method": "Target.createTarget",
		"params": map[string]interface{}{
			"url": url,
		},
	}

	response, err := c.sendCommand(createTargetCmd)
	if err != nil {
		return nil, err
	}

	if response["error"] != nil {
		return nil, parseError(response["error"])
	}

	targetID := response["result"].(map[string]interface{})["targetId"].(string)
	return &Page{ID: targetID}, nil
}

// Navigate navigates a page to the specified URL.
func (c *Client) Navigate(page *Page, url string) error {
	// Send a comand to navigate the page.
	navigateCmd := map[string]interface{}{
		"id":     2,
		"method": "Page.navigate",
		"params": map[string]interface{}{
			"sessionId": page.ID,
			"url":       url,
		},
	}

	response, err := c.sendCommand(navigateCmd)
	if err != nil {
		return err
	}

	if response["error"] != nil {
		return parseError(response["error"])
	}

	return nil
}

// ClosePage closes a page.
func (c *Client) ClosePage(page *Page) error {
	// Send a command to close the page.
	closePageCmd := map[string]interface{}{
		"id":     3,
		"method": "Target.closeTarget",
		"params": map[string]interface{}{
			"targetId": page.ID,
		},
	}

	response, err := c.sendCommand(closePageCmd)
	if err != nil {
		return err
	}

	if response["error"] != nil {
		return parseError(response["error"])
	}

	return nil
}

// ListPages returns a slice of open pages along with additional information.
func (c *Client) ListPages() ([]*Page, error) {
	listPagesCmd := map[string]interface{}{
		"id":     4,
		"method": "Target.getTargets",
	}

	response, err := c.sendCommand(listPagesCmd)
	if err != nil {
		return nil, err
	}

	if response["error"] != nil {
		return nil, parseError(response["error"])
	}

	targetInfos, ok := response["result"].(map[string]interface{})["targetInfos"].([]interface{})
	if !ok {
		return nil, errors.New("failed to parse targetInfos array")
	}

	pages := make([]*Page, 0)

	for _, targetInfo := range targetInfos {
		targetInfoMap, ok := targetInfo.(map[string]interface{})
		if !ok {
			continue
		}

		targetType, ok := targetInfoMap["type"].(string)
		if !ok || targetType != "page" {
			continue
		}

		targetID, ok := targetInfoMap["targetId"].(string)
		if !ok {
			continue
		}

		pageInfo := &Page{
			ID:    targetID,
			Title: targetInfoMap["title"].(string),
			URL:   targetInfoMap["url"].(string),
			Type:  targetInfoMap["type"].(string),
		}

		pages = append(pages, pageInfo)
	}

	return pages, nil
}

// GetPageByID retrieves a page by its ID.
func (c *Client) GetPageByID(pageID string) (*Page, error) {
	pages, err := c.ListPages()
	if err != nil {
		return nil, err
	}

	for _, pageInfo := range pages {
		if pageInfo.ID == pageID {
			return pageInfo, nil
		}
	}

	return nil, errors.New("page not found")
}

func (c *Client) sendCommand(command map[string]interface{}) (map[string]interface{}, error) {
	if err := c.conn.WriteJSON(command); err != nil {
		return nil, err
	}

	var response map[string]interface{}
	if err := c.conn.ReadJSON(&response); err != nil {
		return nil, err
	}

	return response, nil
}

func parseError(errObj interface{}) error {
	errMap, ok := errObj.(map[string]interface{})
	if !ok {
		return nil
	}
	code, _ := errMap["code"].(float64)
	message, _ := errMap["message"].(string)
	return &CommandError{Code: int(code), Message: message}
}

type CommandError struct {
	Code    int
	Message string
}

func (e *CommandError) Error() string {
	return e.Message
}
