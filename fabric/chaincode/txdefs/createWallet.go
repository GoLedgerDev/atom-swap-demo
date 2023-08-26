package txdefs

import (
	"encoding/json"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/goledgerdev/token-cc/chaincode/utils"
)

// POST Method
var CreateWallet = tx.Transaction{
	Tag:         "createWallet",
	Label:       "Create Wallet",
	Description: "Create a new wallet",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required:    true,
			Tag:         "address",
			Label:       "Address",
			Description: "Address of the account",
			DataType:    "string",
		},
		{
			Required:    true,
			Tag:         "label",
			Label:       "Label",
			Description: "Label of the account",
			DataType:    "string",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		address, _ := req["address"].(string)
		label, _ := req["label"].(string)

		// check if address is valid
		_, err := utils.CheckPublicKey(address)
		if err != nil {
			return nil, errors.WrapError(err, "invalid address")
		}

		wallet, err := assets.NewAsset(map[string]interface{}{
			"@assetType":     "wallet",
			"address":        address,
			"label":          label,
			"goTokenBalance": "0",
		})

		if err != nil {
			return nil, errors.WrapError(err, "failed to create wallet")
		}

		walletMap, err := wallet.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to create wallet")
		}

		walletBytes, err := json.Marshal(walletMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to create wallet")
		}

		return walletBytes, nil
	},
}
