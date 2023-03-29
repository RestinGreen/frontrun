// To run the side .go file to test stuff
// contains random function used to test how stuff works
// clear && go build -o test ./cmd/test/ && ./test
package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/RestinGreen/frontrun/pkg/chain"
	"github.com/RestinGreen/frontrun/pkg/flashbots"
	"github.com/RestinGreen/frontrun/pkg/general"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)


func checkbundle() {

	pk, _ := crypto.HexToECDSA("...")
	a := flashbots.NewBundleStat(true)

	a.GetBundleStat("...", 8624385, pk)

}

func nonce() {
	constants := general.NewConstants()
	connection := general.NewConnection(constants.IpcEndpoint)
	victim := chain.NewWallet(connection.EthClient)
	
	nonce, err := connection.EthClient.NonceAt(context.Background(), victim.Address, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(nonce)

}

func a() {
	gen := general.NewConstants()
	conn := general.NewConnection(gen.IpcEndpoint)
	
	nr, err := conn.EthClient.BlockNumber(gen.Ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(nr)
}

func steal() {

	constants := general.NewConstants()
	fmt.Println(constants.Testnet)

	thief_pk, err := crypto.HexToECDSA(os.Getenv("THIEF_PK"))
	write, err := bind.NewKeyedTransactorWithChainID(thief_pk, big.NewInt(5))
	if err != nil {
		fmt.Println("Failed to setup new keyed transaction. Exiting.")
		panic(err)
	}
	
	write.Context = context.Background()
	write.From = constants.Thief
	write.Value = big.NewInt(0)
	write.GasLimit = 100_000
	write.GasPrice = big.NewInt(57_365_442_668)
	
	
	if err != nil {
		fmt.Println("error reading thief pk." )
		panic(err)
	}
	
	conn := general.NewConnection("/root/goerli/execution/chaindata/geth.ipc")
	nr, err := conn.EthClient.BlockNumber(context.Background())
	if err != nil {
		fmt.Println("error getting block number")
		panic(err)
	}
	fmt.Println("block number: ", nr)

	victim := chain.NewWallet(conn.EthClient)
	coin := chain.NewContract(constants, conn.EthClient)
	tx, err := coin.CoinContract.TransferFrom(write, victim.Address, constants.Thief, big.NewInt(220329999999998874))

	fmt.Println("thief tx hash: ", tx.Hash())
}

func main() {

	steal()
}
