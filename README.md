# ghoul - Chrome DevTools Protocol wrapper for Golang

This is a Go client for interacting with the Chrome DevTools Protocol (CDP). It provides methods for creating and managing CDP pages.

## Installation

To use ghoul, you need to have Go installed and set up a Go module for your project. Then, you can add the client as a dependency:

```shell
go get github.com/jonasrdl/ghoul/cdp
```

## Usage
1. Import the necessary packages
```go
import (
	"github.com/gorilla/websocket"
	"github.com/jonasrdl/ghoul/cdp"
)
```
2. Create a WebSocket connection to the Chrome DevTools
```go
wsConn, err := websocket.Dial("ws://localhost:9222/devtools/browser")
if err != nil {
    log.Fatal(err)
}
defer wsConn.Close()
```
3. Create a new CDP client
```go
client := cdp.NewClient(wsConn)
```
4. Use the client to create or close pages
```go
page, err := client.CreatePage("https://www.example.com")
if err != nil {
    log.Fatal(err)
}

err = client.ClosePage(page)
if err != nil {
    log.Fatal(err)
}
```

## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or a pull request