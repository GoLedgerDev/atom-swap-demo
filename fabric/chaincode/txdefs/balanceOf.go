package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
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
			DataType:    "->wallet",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error

		wallet, _ := req["address"].(assets.Key)

		walletObj, err := wallet.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting wallet asset")
		}

		balance, ok := walletObj.GetProp("goTokenBalance").(string)
		if !ok {
			return nil, errors.NewCCError("error getting wallet balance", http.StatusInternalServerError)
		}

		response := make(map[string]interface{})
		response["balance"] = balance
		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling response")
		}

		return responseJSON, nil
	},
}
