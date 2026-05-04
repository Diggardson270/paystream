package main

import "log"

func main() {
	log.Println("paystream-worker starting")
	// Background jobs land here: Horizon polling, signing queue, webhook delivery.
	select {}
}
