package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/goledgerdev/token-cc/chaincode/utils"
)

// GET Method
var BalanceOf = tx.Transaction{
	Tag:         "balanceOf",
	Label:       "Balance Of",
	Description: "BalanceOf returns the token balance of a given account",
	Method:      "GET",

	Args: []tx.Argument{
		{
			Required:    true,
			Tag:         "address",
			Label:       "Address",
			Description: "Address of the account to be checked",
			DataType:    "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error

		address, _ := req["address"].(string)

		// check if address is valid
		_, err = utils.CheckPublicKey(address)
		if err != nil {
			return nil, errors.WrapError(nil, "invalid address")
		}

		wallet, err := assets.NewKey(map[string]interface{}{
			"@assetType": "wallet",
			"address":    address,
		})
		if err != nil {
			return nil, errors.WrapError(err, "error creating wallet key")
		}

		exists, err := wallet.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error checking if wallet asset exists")
		}

		response := map[string]interface{}{
			"balance": "0",
		}

		if exists {
			walletObj, err := wallet.Get(stub)
			if err != nil {
				return nil, errors.WrapError(err, "error getting wallet asset")
			}

			balance, ok := walletObj.GetProp("goTokenBalance").(string)
			if !ok {
				return nil, errors.NewCCError("error getting wallet balance", http.StatusInternalServerError)
			}

			response["balance"] = balance
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling response")
		}

		return responseJSON, nil
	},
}
