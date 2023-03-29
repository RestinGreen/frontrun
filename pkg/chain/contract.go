package chain

import (
	"fmt"
	"math/big"

	"github.com/RestinGreen/frontrun/pkg/binding/usdtmainnet"
	"github.com/RestinGreen/frontrun/pkg/binding/wethgoerli"
	"github.com/RestinGreen/frontrun/pkg/general"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20 interface {
	Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error)
	TransferFrom(opts *bind.TransactOpts, src common.Address, dst common.Address, wad *big.Int) (*types.Transaction, error)
	Transfer(opts *bind.TransactOpts, _to common.Address, _value *big.Int) (*types.Transaction, error)
	BalanceOf(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)
	
}

type Contract struct {
	CoinContract ERC20
}

func NewContract(constants general.Const, ethClient *ethclient.Client) *Contract {

	var c Contract
	var err error

	if constants.Testnet {
		c.CoinContract, err = wethgoerli.NewWethgoerli(constants.GetCoinAddress(), ethClient)
	} else {
		c.CoinContract, err = usdtmainnet.NewUsdtmainnet(constants.GetCoinAddress(), ethClient)
	}
	if err != nil {
		fmt.Println("Failed to create erc20 contract binding. Exiting.")
		panic(err)
	}

	return &c
}
