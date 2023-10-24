package main

import (
	"fmt"
	"github.com/jonasrdl/ghoul/cdp"
	"log"
	"time"
)

func main() {
	client, err := cdp.StartChromiumAndConnect()
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	defer client.Close()

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
