package main

import (
	txdefs "github.com/goledgerdev/token-cc/chaincode/txdefs"

	tx "github.com/goledgerdev/cc-tools/transactions"
)

var txList = []tx.Transaction{
	txdefs.Mint,
}
