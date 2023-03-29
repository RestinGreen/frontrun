package main

import (
	"fmt"

	"github.com/RestinGreen/frontrun/pkg/binding"
	"github.com/RestinGreen/frontrun/pkg/chain"
	"github.com/RestinGreen/frontrun/pkg/general"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

func main() {
	fmt.Println("Hunting for thiefs >.<\n ")

	constants := general.NewConstants()
	connection := general.NewConnection(constants.IpcEndpoint)
	abiHandler := binding.NewAbiHandler(constants.Testnet)
	coin := chain.NewContract(constants, connection.EthClient)
	victim := chain.NewWallet(connection.EthClient)
	interact := chain.NewInteract(victim, connection.EthClient)
	// userstats := flashbots.NewStatistics(constants.Testnet)

	// simulation := flashbots.NewSimulation(constants.Testnet)
	// bundler := flashbots.NewBundler(constants.Testnet)
	// bundleStat := flashbots.NewBundleStat(constants.Testnet)

	// blockNumber, err := connection.EthClient.BlockNumber(constants.Ctx)
	// userstats.GetStats(blockNumber, constants.FlashbotsAuthPK)

	balance, err := coin.CoinContract.BalanceOf(interact.Read, victim.Address)
	if err != nil {
		fmt.Println("Failed to get victim coin balance. Exiting.")
		panic(err)
	}

	fmt.Println("Current balance to save: ", balance)
	victim.Balance = balance
	go victim.KeepUpToDate(coin, interact.Read)

	txs := make(chan *types.Transaction)
	subscriber := gethclient.New(connection.RpcClient)

	_, err = subscriber.SubscribeFullPendingTransactions(constants.Ctx, txs)
	if err != nil {
		fmt.Println("Failed to subsctibe to pending transactions. Exiting.")
		panic(err)
	}

	for tx := range txs {
		var from common.Address
		signer := types.LatestSignerForChainID(tx.ChainId())
		from, err := types.Sender(signer, tx)
		if err != nil {
			fmt.Println(err)
			fmt.Println()
			continue
		}
		if from.Hex() == constants.Thief.Hex() && tx.To().Hex() == constants.GetCoinAddress().Hex() {
			fmt.Println("found thief")
			fmt.Println("hash:\t", tx.Hash())
			fmt.Println("coin:\t", tx.To())
			fmt.Println("tx type: \t", tx.Type())
			fmt.Println("thief:\t", from)
			fmt.Println()

			method, err := abiHandler.Abi.MethodById(tx.Data()[0:4])
			if err != nil {
				fmt.Println("Method not found.")
			}
			if method.Name == "transferFrom" {
				fmt.Println("Method found:\t", method.Name)
				if common.BytesToAddress(tx.Data()[5:36]) == victim.Address {
					fmt.Println("Stealing from me.\nBuilding frontrun.")

					interact.CopyTxData(victim, *tx)
					frontrunTx, err := coin.CoinContract.Transfer(interact.Write, constants.Safe, victim.Balance)
					if err != nil {
						fmt.Println("Error creating frontrun tx.")
						fmt.Println(err)
						continue
					}
					fmt.Println("Ftonrun tx sent: ", frontrunTx.Hash())
					// blockNumber, err := connection.EthClient.BlockNumber(constants.Ctx)

					// signedThiefTx, _ := tx.MarshalBinary()
					// signedFrontrunTx, _ := frontrunTx.MarshalBinary()
					// fmt.Println("t ", hex.EncodeToString(signedThiefTx))

					// for i := 1; i <= 30; i++ {

					// 	go func(i int, tx *types.Transaction) {
					// 		simulation.NewCallRequest(
					// 			[]string{"0x" + hex.EncodeToString(signedFrontrunTx), "0x" + hex.EncodeToString(signedThiefTx)},
					// 			blockNumber+uint64(i),
					// 			constants.FlashbotsAuthPK)
					// 		bundler.SendBundle(
					// 			[]string{"0x" + hex.EncodeToString(signedFrontrunTx), "0x" + hex.EncodeToString(signedThiefTx)},
					// 			blockNumber+uint64(i),
					// 			[]string{tx.Hash().String(), frontrunTx.Hash().String()},
					// 			constants.FlashbotsAuthPK,
					// 		)
					// 		fmt.Println("Target block number: ", blockNumber+uint64(i))
					// 	}(i, tx)
					// }
				}
			}
		}
	}
}
