package main

// This file is mainly in use for testing currently

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jonasrdl/ghoul/cdp"
	"log"
	"time"
)

func main() {
	ctx := context.Background()

	chromiumWsURL := "ws://127.0.0.1:9222/devtools/browser/ae0b92a8-3fdb-48f1-81f4-52c17650239c"

	conn, _, err := websocket.DefaultDialer.DialContext(ctx, chromiumWsURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := cdp.NewClient(conn)

	// Create a new page.
	page, err := client.CreatePage("https://google.de")
	if err != nil {
		log.Fatalf("Error creating page: %v", err)
	}
	log.Printf("Created page with ID: %s\n", page.ID)

	pages, err := client.ListPages()
	if err != nil {
		log.Fatalf("Error listing pages: %v", err)
	}

	fmt.Println("Open Pages:")
	for _, page := range pages {
		fmt.Printf("Page ID: %s\n", page.ID)
		fmt.Printf("Title: %s\n", page.Title)
		fmt.Printf("URL: %s\n", page.URL)
		fmt.Printf("Type: %s\n", page.Type)
		fmt.Println("------------------------------------")
	}

	// Wait for 3 seconds.
	time.Sleep(10 * time.Second)

	// Close the page.
	err = client.ClosePage(page)
	if err != nil {
		log.Fatalf("Error closing page: %v", err)
	}
	fmt.Printf("Closed page with ID: %s\n", page.ID)
}
