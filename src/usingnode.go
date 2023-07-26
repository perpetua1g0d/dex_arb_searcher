package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

func runNodeProvider() {
	const contractAddress = "0x3416cf6c708da44db2624d63ea0aaef7113527c6"
	const rpc = "https://mainnet.infura.io/v3/2554b87d4a874952ae9ce3dcdc9b5fb4"

	fdata, err := ioutil.ReadFile("abi.json")
	HandleError(err)

	contractAbi, err := abi.JSON(bytes.NewReader(fdata))
	HandleError(err)

	// for key, value := range contractMethods {
	// 	fmt.Printf("Found method: %s:%s\n", key, value)
	// }

	// contractMethods := contractAbi.Methods
	// fmt.Println(contractAbi.MethodById(contractMethods["fee"].ID))

	client, err := ethclient.Dial(rpc)
	HandleError(err)

	packed, err := contractAbi.Pack("slot0")
	HandleError(err)
	address := common.HexToAddress(contractAddress)

	callMsg := ethereum.CallMsg{
		To:   &address,
		Data: packed,
	}

	response, err := client.CallContract(context.Background(), callMsg, nil)
	HandleError(err)

	fmt.Printf("response: %v\n", response)
	fmt.Printf("hexutil.Encode(response): %v\n", hexutil.Encode(response))
	
}
