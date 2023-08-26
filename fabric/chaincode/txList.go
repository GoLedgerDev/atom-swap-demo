package main

import (
	txdefs "github.com/goledgerdev/token-cc/chaincode/txdefs"

	tx "github.com/goledgerdev/cc-tools/transactions"
)

var txList = []tx.Transaction{
	txdefs.CreateWallet,
	txdefs.Mint,
	txdefs.Burn,
	txdefs.Transfer,
	txdefs.TotalSupply,
	txdefs.BalanceOf,

	// TODO: implement these transactions
	// txdefs.Allowance,
	// txdefs.Approve,
	// txdefs.IncreaseAllowance,
	// txdefs.DecreaseAllowance,
	// txdefs.TransferFrom,
}
