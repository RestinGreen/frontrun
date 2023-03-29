package general

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
)

// Constant variables used in by the bot.
// Values come from .env file, which is in the root of the project.
type Const struct {
	FlashbotsAuthPK    *ecdsa.PrivateKey
	Thief              common.Address
	Safe               common.Address
	testnetCoinAddress common.Address
	mainnetCoinAddress common.Address
	Testnet            bool

	IpcEndpoint    string
	Ctx            context.Context
	TransferFromID []byte
}

func NewConstants() Const {

	var c Const
	var err error

	test := os.Getenv("TESTNET")
	if strings.Compare(test, "false") == 0 {
		c.Testnet = false
		err = godotenv.Load("./.env")
	} else {
		c.Testnet = true
		err = godotenv.Load("./.env.test")
	}


	if err != nil {
		fmt.Println("Couldn't load .env file. Exiting.\n", err)
		os.Exit(1)
	}

	
	c.Thief = common.HexToAddress(os.Getenv("THIEF_WALLET"))
	c.FlashbotsAuthPK, err = crypto.HexToECDSA(os.Getenv("FLASHBOTS_AUTH_PK"))
	if err != nil {
		fmt.Println("Error loadin Flashbots Authorization PK. Exiting.")
		panic(err)
	}

	c.Safe = common.HexToAddress(os.Getenv("SAFE_WALLET"))
	c.testnetCoinAddress = common.HexToAddress(os.Getenv("TESTNET_COIN_ADDRESS"))
	c.mainnetCoinAddress = common.HexToAddress(os.Getenv("MAINNET_COIN_ADDRESS"))

	c.IpcEndpoint = os.Getenv("IPC_ENDPOINT")
	c.Ctx = context.Background()
	c.TransferFromID = []byte{35, 184, 114, 221}

	

	return c
}

func (c Const) GetCoinAddress() common.Address {

	if c.Testnet {
		return c.testnetCoinAddress
	} else {
		return c.mainnetCoinAddress
	}
}
