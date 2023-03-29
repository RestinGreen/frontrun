package chain

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Wallet struct {
	Address    common.Address
	PrivateKey *ecdsa.PrivateKey
	Balance    *big.Int
	nonce      *big.Int

	nonceMut sync.Mutex
	balMut   sync.Mutex

	ethClient *ethclient.Client
}

func NewWallet(ethClient *ethclient.Client) *Wallet {

	var w Wallet
	var err error

	w.ethClient = ethClient
	w.Address = common.HexToAddress(os.Getenv("VICTIM_WALLET"))
	w.PrivateKey, err = crypto.HexToECDSA(os.Getenv("VICTIM_PK"))
	if err != nil {
		fmt.Println("Error loadin Victim PK. Exiting.")
		panic(err)
	}

	return &w

}

func (w *Wallet) GetBalance() *big.Int {
	w.balMut.Lock()
	defer w.balMut.Unlock()
	return w.Balance
}

func (w *Wallet) UpdateBalance(coin *Contract, read *bind.CallOpts) {

	newBal, err := coin.CoinContract.BalanceOf(read, w.Address)
	if err != nil {
		fmt.Println("Error getting balance from chain.\n", err)
	}
	if w.Balance.Cmp(newBal) != 0 {
		w.balMut.Lock()
		w.Balance = newBal
		fmt.Println("New balance to save: ", newBal)
		w.balMut.Unlock()
	}
	
}

func (w *Wallet) GetNonce() *big.Int {

	w.nonceMut.Lock()
	defer w.nonceMut.Unlock()
	return w.nonce
}

func (w *Wallet) UpdateNonce() {
	nonce, err := w.ethClient.PendingNonceAt(context.Background(), w.Address)
	if err != nil {
		fmt.Println("Failed to querry new wallet nonce. ", err)
		os.Exit(1)
	}
	w.nonceMut.Lock()
	w.nonce = big.NewInt(int64(nonce))
	w.nonceMut.Unlock()
	
}

func (w *Wallet) KeepUpToDate(coin *Contract, read *bind.CallOpts) {
	newBlocks := make(chan *types.Header)
	w.ethClient.SubscribeNewHead(context.Background(), newBlocks)

	for newBlock := range newBlocks {
		fmt.Println("new block: ", newBlock.Number)
		w.UpdateNonce()
		w.UpdateBalance(coin, read)
	}
}
