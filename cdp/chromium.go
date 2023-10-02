package cdp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

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

type CommandError struct {
	Code    int
	Message string
}

func (e *CommandError) Error() string {
	return e.Message
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
