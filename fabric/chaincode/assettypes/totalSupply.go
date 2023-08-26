package assettypes

import (
	"github.com/goledgerdev/cc-tools/assets"
)

const TotalSupplyKey = "totalSupply"

var TotalSupply = assets.AssetType{
	Tag:         "totalSupply",
	Label:       "Total Supply",
	Description: "Total supply of an ERC20 token",

	Props: []assets.AssetProp{
		{
			// Primary key
			Required:     true,
			IsKey:        true,
			Tag:          "totalSupplyKey",
			Label:        "Total Supply Key",
			DataType:     "string",
			DefaultValue: "totalSupply",
		},
		{
			Required: true,
			Tag:      "supply",
			Label:    "Supply",
			DataType: "string",
		},
	},
}
