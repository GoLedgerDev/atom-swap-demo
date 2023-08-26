package txdefs

import (
	"encoding/json"
	"strconv"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/goledgerdev/token-cc/chaincode/utils"
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
			DataType:    "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		amountStr, _ := req["amount"].(string)
		from, _ := req["from"].(string)

		// check if amount is valid integer
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return nil, errors.WrapError(err, "invalid amount, send a positive integer")
		}

		if amount <= 0 {
			return nil, errors.WrapError(nil, "invalid amount, send a positive integer")
		}

		// check if address is valid
		_, err = utils.CheckPublicKey(from)
		if err != nil {
			return nil, errors.WrapError(nil, "invalid address")
		}

		// retrieve wallet
		wallet, err := assets.NewKey(map[string]interface{}{
			"@assetType": "wallet",
			"address":    from,
		})

		// check if wallet exists
		exists, err := wallet.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error checking if wallet exists")
		}

		currentAmount := 0
		if exists {
			// get current amount
			walletAsset, err := wallet.Get(stub)
			if err != nil {
				return nil, err
			}

			c, nerr := strconv.Atoi(walletAsset.GetProp("goTokenBalance").(string))
			if nerr != nil {
				return nil, errors.WrapError(err, "error converting current amount to integer")
			}

			currentAmount = c
		}

		newAmount := currentAmount - amount
		if newAmount < 0 {
			return nil, errors.WrapError(nil, "insufficient funds")
		}

		walletMap := map[string]interface{}{
			"@assetType":     "wallet",
			"address":        from,
			"goTokenBalance": strconv.Itoa(newAmount),
		}

		w, err := assets.NewAsset(walletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error creating wallet asset")
		}

		walletInLedger, err := w.Put(stub)
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
