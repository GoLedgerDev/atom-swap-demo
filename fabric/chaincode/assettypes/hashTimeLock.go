package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

var HashTimeLock = assets.AssetType{
	Tag:         "hashTimeLock",
	Label:       "Hash Time Lock",
	Description: "Hash Time Lock",

	Props: []assets.AssetProp{
		{
			// Primary key
			IsKey:    true,
			Required: true,
			Tag:      "id",
			Label:    "Lock ID",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "fromWallet",
			Label:    "From Wallet",
			DataType: "->wallet",
		},
		{
			Required: true,
			Tag:      "toWallet",
			Label:    "To Wallet",
			DataType: "->wallet",
		},
		{
			Required: true,
			Tag:      "amount",
			Label:    "Amount",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "hashlock",
			Label:    "Hash Lock",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "timelock",
			Label:    "Time Lock",
			DataType: "datetime",
		},
		{
			Tag:      "secret",
			Label:    "Hash Secret",
			DataType: "string",
		},
	},
}
