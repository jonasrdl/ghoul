package cdp

import "github.com/gorilla/websocket"

// Page represents a CDP page.
type Page struct {
	ID string `json:"id"`
}

// Client provides methods to interact with the CDP.
type Client struct {
	conn *websocket.Conn
}

// NewClient creates a new CDP client.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn}
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
