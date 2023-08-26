package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

var Wallet = assets.AssetType{
	Tag:         "wallet",
	Label:       "Wallet",
	Description: "A wallet will hold the balance of an ERC20 token",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required: true,
			IsKey:    true,
			Tag:      "address",
			Label:    "Address",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "goTokenBalance",
			Label:    "GoToken Balance",
			DataType: "string",
		},
	},
}
