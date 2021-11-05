package http

import "fmt"

const (
	httpAuthHeader      = "Authorization"
	httpAuthValueFormat = "TermlyV1, PublicKey=%s, Signature=%s"
)

func createAuthValue(publicKey, signature string) string {
	return fmt.Sprintf(httpAuthValueFormat, publicKey, signature)
}
