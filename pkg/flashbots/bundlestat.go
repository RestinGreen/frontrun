package flashbots

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type BundleStatRequest struct {
	Jsonrpc string             `json:"jsonrpc"`
	ID      int                `json:"id"`
	Method  string             `json:"method"`
	Params  []BundleStatParams `json:"params"`
}

type BundleStatParams struct {
	BundleHash  string `json:"bundleHash"`
	BlockNumber string `json:"blockNumber"`
}

type Builder struct {
	Pubkey    string `json:"pubkey"`
	Timestamp string `json:"timestamp"`
}

type BundleStatResponse struct {
	IsHighPriority         bool      `json:"isHighPriority"`
	IsSimulated            bool      `json:"isSimulated"`
	SimulatedAt            time.Time `json:"simulatedAt"`
	ReceivedAt             time.Time `json:"receivedAt"`
	ConsideredByBuildersAt []Builder `json:"consideredByBuildersAt"`
	SealedByBuildersAt     []Builder `json:"sealedByBuildersAt"`
}

type BundleStat struct {
	request  BundleStatRequest
	response BundleStatResponse

	url string
}

func NewBundleStat(isTestnet bool) *BundleStat {
	var b BundleStat

	if isTestnet {
		b.url = "https://relay.flashbots.net"
	} else {
		b.url = "https://relay-goerli.flashbots.net"
	}

	b.request = BundleStatRequest{
		Jsonrpc: "2.0",
		ID:      1,
		Method:  "flashbots_getBundleStatsV2",
		Params:  []BundleStatParams{},
	}

	return &b
}

func (b *BundleStat) GetBundleStat(bundleHash string, blockNumber uint64, flashbotsAuthPk *ecdsa.PrivateKey) {
	
	params := BundleStatParams {
		BundleHash: bundleHash,
		BlockNumber: fmt.Sprintf("0x%x", blockNumber),
	}
	b.request.Params = append(b.request.Params, params)

	bodyBytes, err := json.Marshal(b.request)
    if err != nil {
        fmt.Println("Failed to marshall bundle stat request struct to json.")
        fmt.Println(err)
    }
    req, err := http.NewRequest("POST", b.url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        fmt.Println("Failed to create bundle stat request.")
        fmt.Println(err)
    }

    flashbotsHeader := CreateHeaderSignature(flashbotsAuthPk, string(bodyBytes))

    req.Header.Set("X-Flashbots-Signature", flashbotsHeader)
	req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to do the bundle stat http request to flashbots.")
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	var res BundleStatResponse
	err = json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println("IsHighPriority:", res.IsHighPriority)
	fmt.Println("IsSimulated:", res.IsSimulated)
	fmt.Println("SimulatedAt:", res.SimulatedAt)
	fmt.Println("ReceivedAt:", res.ReceivedAt)

	fmt.Println("ConsideredByBuildersAt:")
	for _, builder := range res.ConsideredByBuildersAt {
		fmt.Println("  Pubkey:", builder.Pubkey)
		fmt.Println("  Timestamp:", builder.Timestamp)
	}

	fmt.Println("SealedByBuildersAt:")
	for _, builder := range res.SealedByBuildersAt {
		fmt.Println("  Pubkey:", builder.Pubkey)
		fmt.Println("  Timestamp:", builder.Timestamp)
	}

}
