package main

import (
	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/token-cc/chaincode/assettypes"
)

var assetTypeList = []assets.AssetType{
	assettypes.Wallet,
	assettypes.TotalSupply,
}
