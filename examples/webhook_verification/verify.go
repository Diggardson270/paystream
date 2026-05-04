package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// VerifyWebhook checks that a PayStream webhook was signed by your endpoint secret.
func VerifyWebhook(payload []byte, signature, secret string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func main() {
	ok := VerifyWebhook([]byte(`{"event":"payout.settled"}`), "sha256=...", "your-secret")
	fmt.Println("valid:", ok)
}
