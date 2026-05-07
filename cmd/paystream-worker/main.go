package main

import (
	"log"
	"time"
)

func main() {
	log.Printf("paystream-worker starting at %s", time.Now().UTC().Format(time.RFC3339))
	// Background jobs land here: Horizon polling, signing queue, webhook delivery.
	select {}
}
