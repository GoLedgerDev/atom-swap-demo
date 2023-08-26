package txdefs

import (
	"encoding/json"
	"strconv"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// POST Method
var Burn = tx.Transaction{
	Tag:         "burn",
	Label:       "Burn",
	Description: "Burn redeems tokens the minter's account balance",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required:    true,
			Tag:         "amount",
			Label:       "Amount",
			Description: "Amount of tokens to be burned",
			DataType:    "string", // Receive as string to avoid precision loss
		},
		{
			Required:    true,
			Tag:         "from",
			Label:       "From",
			Description: "Address of the account that will burn the tokens",
			DataType:    "->wallet",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		amountStr, _ := req["amount"].(string)
		wallet, _ := req["from"].(assets.Key)

		// check if amount is valid integer
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return nil, errors.WrapError(err, "invalid amount, send a positive integer")
		}

		if amount <= 0 {
			return nil, errors.WrapError(nil, "invalid amount, send a positive integer")
		}

		// get current amount
		walletAsset, err := wallet.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting wallet asset")
		}

		currentAmount, nerr := strconv.Atoi(walletAsset.GetProp("goTokenBalance").(string))
		if nerr != nil {
			return nil, errors.WrapError(err, "error converting current amount to integer")
		}

		newAmount := currentAmount - amount
		if newAmount < 0 {
			return nil, errors.WrapError(nil, "insufficient funds")
		}

		walletMap := map[string]interface{}{
			"goTokenBalance": strconv.Itoa(newAmount),
		}

		walletInLedger, err := wallet.Update(stub, walletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error putting wallet asset")
		}

		walletJSON, err := json.Marshal(walletInLedger)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling wallet asset")
		}

		// remove tokens from total supply
		err = AddAmountToTotalSupply(stub, -amount)
		if err != nil {
			return nil, errors.WrapError(err, "error removing amount from total supply")
		}

		return walletJSON, nil
	},
}
