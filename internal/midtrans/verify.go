package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"os"
)

func VerifySignature(p WebhookPayload) bool {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return false
	}

	raw := p.OrderID +
		p.StatusCode +
		p.GrossAmount +
		serverKey

	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])

	return expected == p.SignatureKey
}
