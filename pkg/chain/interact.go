package chain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Interact struct {
	Write *bind.TransactOpts
	Read  *bind.CallOpts
}

func NewInteract(wallet *Wallet, ethClient *ethclient.Client) *Interact {

	var i Interact

	chainId, err := ethClient.ChainID(context.Background())
	if err != nil {
		fmt.Println("Failed to query chain id.  Exiting.")
		panic(err)
	}

	i.Read = &bind.CallOpts{
		Pending: false,
		From:    wallet.Address,
	}

	i.Write, err = bind.NewKeyedTransactorWithChainID(wallet.PrivateKey, chainId)
	if err != nil {
		fmt.Println("Failed to setup new keyed transaction. Exiting.")
		panic(err)
	}

	i.Write.Context = context.Background()
	i.Write.From = wallet.Address
	i.Write.Value = big.NewInt(0)
	i.Write.GasLimit = 100_000
	// i.Write.NoSend = false

	return &i
}

func (i *Interact) CopyTxData(wallet *Wallet, thiefTx types.Transaction) {

	i.Write.Nonce = wallet.GetNonce()
	if thiefTx.Type() == 2 {
		
		i.Write.GasFeeCap = big.NewInt(0)
		i.Write.GasTipCap = big.NewInt(0)
		i.Write.GasFeeCap.Add(thiefTx.GasFeeCap(), big.NewInt(1))
		i.Write.GasTipCap.Mul(thiefTx.GasTipCap(), big.NewInt(2))
	} else {
		i.Write.GasPrice = big.NewInt(0)
		i.Write.GasPrice.Add(thiefTx.GasPrice(), big.NewInt(1_054_924_419)) //add ~1 gwei
	}

}

func (i *Interact) SetSend(send bool) {
	i.Write.NoSend = send
}
