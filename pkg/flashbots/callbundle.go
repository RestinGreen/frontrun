package flashbots

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type CallBundleRequest struct {
	JSONRPC string             `json:"jsonrpc"`
	ID      int                `json:"id"`
	Method  string             `json:"method"`
	Params  []CallBundleParams `json:"params"`
}

type CallBundleParams struct {
	Txs              []string `json:"txs"`
	BlockNumber      string   `json:"blockNumber"`
	StateBlockNumber string   `json:"stateBlockNumber"`
	Timestamp        int64    `json:"timestamp"`
}

type CallBundleResponse struct {
	JSONRPC string       `json:"jsonrpc"`
	ID      int       `json:"id"`
	Result  CallBundleResult `json:"result"`
}

type CallBundleResult struct {
	BundleGasPrice    string               `json:"bundleGasPrice"`
	BundleHash        string               `json:"bundleHash"`
	CoinbaseDiff      string               `json:"coinbaseDiff"`
	EthSentToCoinbase string               `json:"ethSentToCoinbase"`
	GasFees           string               `json:"gasFees"`
	Results           []CallBundleResultTx `json:"results"`
	StateBlockNumber  int64                `json:"stateBlockNumber"`
	TotalGasUsed      int64                `json:"totalGasUsed"`
}

type CallBundleResultTx struct {
	CoinbaseDiff      string `json:"coinbaseDiff"`
	EthSentToCoinbase string `json:"ethSentToCoinbase"`
	FromAddress       string `json:"fromAddress"`
	GasFees           string `json:"gasFees"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           int64  `json:"gasUsed"`
	ToAddress         string `json:"toAddress"`
	TxHash            string `json:"txHash"`
	Value             string `json:"value"`
}

type Simulation struct {
	request  CallBundleRequest
	response CallBundleResponse

	url string
}

func NewSimulation(isTestnet bool) *Simulation {

	var sim Simulation
	if isTestnet {
		sim.url = "https://relay-goerli.flashbots.net"
		} else {
		sim.url = "https://relay.flashbots.net"
	}

	sim.request = CallBundleRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "eth_callBundle",
		Params:  []CallBundleParams{},
	}


	return &sim
}

func (s *Simulation) NewCallRequest(txs []string, blockNumber uint64, flashbotsAuthPk *ecdsa.PrivateKey) {

	params := CallBundleParams{
		Txs:              txs,
		BlockNumber:      fmt.Sprintf("0x%x", blockNumber),
		StateBlockNumber: "latest",
		Timestamp:        time.Now().Unix(),
	}

	s.request.Params = [] CallBundleParams{}
	s.request.Params = append(s.request.Params, params)

	bodyBytes, err := json.Marshal(s.request)
	if err != nil {
		fmt.Println("Failed to marshall simulation request struct to json.")
		fmt.Println(err)
	}
	req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(bodyBytes))

	if err != nil {
		fmt.Println("Faield to create simulation request.")
		fmt.Println(err)
	}

	flashbotsHeader := CreateHeaderSignature(flashbotsAuthPk, string(bodyBytes))

	req.Header.Set("X-Flashbots-Signature", flashbotsHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to do the simulation http request to flashbots.")
		fmt.Println(err)

	}
	defer resp.Body.Close()

	var res CallBundleResponse

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		fmt.Println("Failed to unmarshall call bundle http response")
		fmt.Println(err)
	}

	fmt.Println("Bundle gas price:", res.Result.BundleGasPrice)
	fmt.Println("Bundle hash:", res.Result.BundleHash)
	fmt.Println("Coinbase diff:", res.Result.CoinbaseDiff)
	fmt.Println("ETH sent to Coinbase:", res.Result.EthSentToCoinbase)
	fmt.Println("Gas fees:", res.Result.GasFees)
	fmt.Println("State block number:", res.Result.StateBlockNumber)
	fmt.Println("Total gas used:", res.Result.TotalGasUsed)

	// print the results
	fmt.Println("Results:")
	for _, result := range res.Result.Results {
		fmt.Println("  Tx hash:", result.TxHash)
		fmt.Println("  From address:", result.FromAddress)
		fmt.Println("  To address:", result.ToAddress)
		fmt.Println("  Value:", result.Value)
		fmt.Println("  Gas price:", result.GasPrice)
		fmt.Println("  Gas used:", result.GasUsed)
		fmt.Println("  Gas fees:", result.GasFees)
		fmt.Println("  Coinbase diff:", result.CoinbaseDiff)
		fmt.Println("  ETH sent to Coinbase:", result.EthSentToCoinbase)
		fmt.Println()
	}

}
