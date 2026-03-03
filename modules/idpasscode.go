package modules

import (
	"crypto/rand"
	"fmt"
	"time"
)

// ======================
// CONSTANT
// ======================
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// ======================
// STRUCT
// ======================
type Order struct {
	ID       string `json:"id"`
	Passcode string `json:"passcode"`
	Name     string `json:"name"`
	Contact  string `json:"contact"`
	Service  string `json:"service"`
	Message  string `json:"message"`
	Status   string `json:"status"`
}

// ======================
// GENERATORS
// ======================

// ORD-YYYYMMDD-XXXXXX
func GenerateOrderID() string {
	return fmt.Sprintf(
		"ORD-%s-%s",
		time.Now().Format("20060102"),
		randomString(6),
	)
}

// 8 char a-zA-Z0-9
func GeneratePasscode() string {
	return randomString(8)
}

// ======================
// INTERNAL
// ======================
func randomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
