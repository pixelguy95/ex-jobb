package main

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcutil"

	"./customtransactions"
	"./rpcutils"
	"github.com/btcsuite/btcd/rpcclient"
)

var connCfg = &rpcclient.ConnConfig{
	Host:         "localhost:18332",
	HTTPPostMode: true,
	DisableTLS:   true,
	User:         "pi",
	Pass:         "kebab",
}

func main() {

	ntfnHandlers := rpcclient.NotificationHandlers{
		OnClientConnected: func() {
			fmt.Println("Connected")
		},
	}

	client, error := rpcclient.New(connCfg, &ntfnHandlers)
	if error != nil {
		fmt.Println(error)
	}

	client.Connect(1)
	clientWraper := rpcutils.New(client)

	clientWraper.UnloadAllWallets()
	clientWraper.LoadWallet("cj_wallet")

	receiver, _ := btcutil.DecodeAddress("mymoiBMCFzTp51fAnFftXrAZqAAhDVMLKd", &chaincfg.TestNet3Params)
	contractDetails, _ := customtransactions.GenerateAtomicSwapContract(receiver, 100000, client)

	stuff, _ := txscript.ExtractAtomicSwapDataPushes(0, contractDetails.Contract)
	recipient, _ := btcutil.NewAddressPubKeyHash(stuff.RecipientHash160[:], &chaincfg.TestNet3Params)
	refund, _ := btcutil.NewAddressPubKeyHash(stuff.RefundHash160[:], &chaincfg.TestNet3Params)
	fmt.Printf("\n\t===Contract details===\nTimelock: \t\t%d\nRecipient:\t\t%s\nRefund address:\t\t%s\nSecret hash:\t\t%s\nSecret bytelength:\t%d\n\n",
		stuff.LockTime,
		recipient.EncodeAddress(),
		refund.EncodeAddress(),
		hex.EncodeToString(stuff.SecretHash[:]),
		stuff.SecretSize)

	fmt.Printf("Secret hex: %x\n\n", contractDetails.Secret)

	fmt.Println("Contract:")
	fmt.Printf("%x\n\n", contractDetails.Contract)

	fmt.Println("Contract transaction:")
	buf := new(bytes.Buffer)
	contractDetails.ContractTx.SerializeNoWitness(buf)
	fmt.Printf("%x\n\n", buf)

	fmt.Println("Refund transaction:")
	buf = new(bytes.Buffer)
	contractDetails.Refund.SerializeNoWitness(buf)
	fmt.Printf("%x\n\n", buf)

	client.Disconnect()
}