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
var Mint = tx.Transaction{
	Tag:         "mint",
	Label:       "Mint",
	Description: "Mint creates new tokens and adds them to minter's account balance",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required:    true,
			Tag:         "amount",
			Label:       "Amount",
			Description: "Amount of tokens to be minted",
			DataType:    "string", // Receive as string to avoid precision loss
		},
		{
			Required:    true,
			Tag:         "to",
			Label:       "To",
			Description: "Address of the account that will receive the minted tokens",
			DataType:    "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		amountStr, _ := req["amount"].(string)
		to, _ := req["to"].(string)

		// check if amount is valid integer
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return nil, errors.WrapError(err, "invalid amount, send an integer")
		}

		// check if address is valid
		_, err = utils.CheckPublicKey(to)
		if err != nil {
			return nil, errors.WrapError(nil, "invalid address")
		}

		// retrieve wallet
		wallet, err := assets.NewKey(map[string]interface{}{
			"@assetType": "wallet",
			"address":    to,
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

		walletMap := map[string]interface{}{
			"@assetType":     "wallet",
			"address":        to,
			"goTokenBalance": strconv.Itoa(currentAmount + amount),
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

		return walletJSON, nil
	},
}
