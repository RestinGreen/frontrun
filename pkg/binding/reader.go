package binding

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type AbiHandler struct {
	Abi       abi.ABI
}

func NewAbiHandler(isTestnet bool) *AbiHandler {

	var err error

	var fileName string
	if isTestnet {
		fileName = "./pkg/binding/wethgoerli/wethgoerli.json"
	} else {
		fileName = "./pkg/binding/usdtmainnet/usdtmainnet.json"
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file. Exiting.")
		panic(err)
	}
	defer file.Close()


	buffer := readFile(file)

	abiObj, err := abi.JSON(bytes.NewReader(buffer))
	if err != nil {
		fmt.Println("Failed to read abi. Exiting.")
		panic(err)
	}
	return &AbiHandler{
		Abi: abiObj,
	}
}

func readFile(file *os.File) []byte {
	buffer := make([]byte, 1024)
    var result []byte
    for {
        n, err := file.Read(buffer)
        if err != nil && err != io.EOF {
            fmt.Println("Error reading file:", err)
            panic(err)
        }
        if n == 0 {
            break
        }
        result = append(result, buffer[:n]...)
    }
	return result
}