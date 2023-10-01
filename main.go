package main

import (
	"context"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jonasrdl/ghoul/cdp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	devToolsURL := "ws://127.0.0.1:9222/devtools/browser/43acd9dc-d954-47a0-884b-678f82df26f0"
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, devToolsURL, nil)
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

	// Navigate to a URL.
	/* err = client.Navigate(page, "https://google.de")
	if err != nil {
		log.Fatalf("Error navigating: %v", err)
	} */

	// Wait for 3 seconds.
	time.Sleep(3 * time.Second)

	// Close the page.
	err = client.ClosePage(page)
	if err != nil {
		log.Fatalf("Error closing page: %v", err)
	}
}
