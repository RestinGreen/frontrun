package flashbots

import (
	"crypto/ecdsa"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Auth struct {
}

func CreateHeaderSignature(pk *ecdsa.PrivateKey, body string) string{
	hashedBody := crypto.Keccak256Hash([]byte(body)).Hex()
	sig, err := crypto.Sign(accounts.TextHash([]byte(hashedBody)), pk)
	if err != nil {
		fmt.Println("Failed to build header signature")
	}
	return crypto.PubkeyToAddress(pk.PublicKey).Hex() + ":" + hexutil.Encode(sig)
}
