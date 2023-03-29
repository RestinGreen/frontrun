package flashbots

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserStatsRequest struct {
	JSONRPC string            `json:"jsonrpc"`
	ID      int               `json:"id"`
	Method  string            `json:"method"`
	Params  []UserStatsParams `json:"params"`
}

type UserStatsParams struct {
	BlockNumber string `json:"blockNumber"`
}

type UserStatsResponse struct {
	IsHighPriority           bool   `json:"isHighPriority"`
	AllTimeValidatorPayments string `json:"allTimeValidatorPayments"`
	AllTimeGasSimulated      string `json:"allTimeGasSimulated"`
	Last7dValidatorPayments  string `json:"last7dValidatorPayments"`
	Last7dGasSimulated       string `json:"last7dGasSimulated"`
	Last1dValidatorPayments  string `json:"last1dValidatorPayments"`
	Last1dGasSimulated       string `json:"last1dGasSimulated"`
}

type Statistics struct {
	request  UserStatsRequest
	response UserStatsResponse

	url string
}

func NewStatistics(isTestnet bool) *Statistics {

	var s Statistics

	s.request.JSONRPC = "2.0"
	s.request.ID = 1
	s.request.Method = "flashbots_getUserStatsV2"
	s.request.Params = []UserStatsParams{}

	if isTestnet {
		s.url = "https://relay.flashbots.net"
	} else {
		s.url = "https://relay-goerli.flashbots.net"
	}

	return &s
}

func (s *Statistics) GetStats(blockNumber uint64, flashbotsAuthPk *ecdsa.PrivateKey) {

	params := UserStatsParams{
		BlockNumber: fmt.Sprintf("0x%x", blockNumber),
	}
	s.request.Params = [] UserStatsParams{}
	s.request.Params = append(s.request.Params, params)

	bodyBytes, err := json.Marshal(s.request)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Failed to marshall request struct to json")
	}
	req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		fmt.Println(err)
		fmt.Println("Faield to create user stats request.")
	}

	flashbotsHeader := CreateHeaderSignature(flashbotsAuthPk, string(bodyBytes))

	req.Header.Set("X-Flashbots-Signature", flashbotsHeader)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to do the user stats http request to flashbots.")
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

    var res UserStatsResponse
    err = json.NewDecoder(resp.Body).Decode(&res)
    if err != nil {
		fmt.Println(err)
        fmt.Println("Failed to unmarshall user stats http response.")
    }

    fmt.Printf("IsHighPriority: %v\n", res.IsHighPriority)
    fmt.Printf("AllTimeValidatorPayments: %s\n", res.AllTimeValidatorPayments)
    fmt.Printf("AllTimeGasSimulated: %s\n", res.AllTimeGasSimulated)
    fmt.Printf("Last7dValidatorPayments: %s\n", res.Last7dValidatorPayments)
    fmt.Printf("Last7dGasSimulated: %s\n", res.Last7dGasSimulated)
    fmt.Printf("Last1dValidatorPayments: %s\n", res.Last1dValidatorPayments)
    fmt.Printf("Last1dGasSimulated: %s\n", res.Last1dGasSimulated)

}
