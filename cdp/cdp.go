package cdp

import (
	"github.com/gorilla/websocket"
)

// Client provides methods to interact with the CDP.
type Client struct {
	conn *websocket.Conn
}

// NewClient creates a new CDP client.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{conn}
}

func (c *Client) Close() error {
	return c.conn.Close()
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
