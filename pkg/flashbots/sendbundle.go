package flashbots

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
)

type SendBundleRequest struct {
	JSONRPC string             `json:"jsonrpc"`
	ID      int                `json:"id"`
	Method  string             `json:"method"`
	Params  []SendBundleParams `json:"params"`
}

type SendBundleParams struct {
	Txs               []string `json:"txs"`
	BlockNumber       string   `json:"blockNumber"`
	MinTimestamp      *int64    `json:"minTimestamp,omitempty"`
	MaxTimestamp      *int64    `json:"maxTimestamp,omitempty"`
	RevertingTxHashes []string `json:"revertingTxHashes"`
	ReplacementUuid   *string   `json:"replacementUuid,omitempty"`
}

type SendBundleResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  SendBundleResult  `json:"result"`
}

type SendBundleResult struct {
	BundleHash string `json:"bundleHash"`
}

type Bundler struct {
    request  SendBundleRequest
    response SendBundleResponse

    url string
}

func NewBundler(isTestnet bool) *Bundler{

    var b Bundler

    if isTestnet {
		b.url = "https://relay.flashbots.net"
	} else {
		b.url = "https://relay-goerli.flashbots.net"
	}


    b.request = SendBundleRequest{
        JSONRPC: "2.0",
        ID: 1,
        Method: "eth_sendBundle",
        Params: []SendBundleParams{},
    }

    return &b
}

func (b *Bundler) SendBundle(txs []string, blockNumber uint64, revertTx []string, flashbotsAuthPk *ecdsa.PrivateKey) string {

    params := SendBundleParams {
        Txs: txs,
        BlockNumber: fmt.Sprintf("0x%x", blockNumber),
        RevertingTxHashes: revertTx,
    }
    b.request.Params = []SendBundleParams{}
    b.request.Params = append(b.request.Params, params)

    bodyBytes, err := json.Marshal(b.request)
    if err != nil {
        fmt.Println("Failed to marshall send bundle request struct to json.")
        fmt.Println(err)
    }
    req, err := http.NewRequest("POST", b.url, bytes.NewBuffer(bodyBytes))
    if err != nil {
        fmt.Println("Failed to create send bundle request.")
        fmt.Println(err)
    }

    flashbotsHeader := CreateHeaderSignature(flashbotsAuthPk, string(bodyBytes))

    req.Header.Set("X-Flashbots-Signature", flashbotsHeader)
	req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to do the send bundle http request to flashbots.")
		fmt.Println(err)
	}
	defer resp.Body.Close()

    var res SendBundleResponse
    err = json.NewDecoder(resp.Body).Decode(&res)

    fmt.Println("Bundlehash: ", res.Result.BundleHash)
    return res.Result.BundleHash
}