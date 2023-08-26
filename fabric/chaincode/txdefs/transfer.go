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
var Transfer = tx.Transaction{
	Tag:         "transfer",
	Label:       "Transfer",
	Description: "Transfer transfers tokens from sender account to recipient account",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required:    true,
			Tag:         "amount",
			Label:       "Amount",
			Description: "Amount of tokens to be transferred",
			DataType:    "string", // Receive as string to avoid precision loss
		},
		{
			Required:    true,
			Tag:         "to",
			Label:       "To",
			Description: "Address of the account that will receive the transferred tokens",
			DataType:    "->wallet",
		},
		{
			Required:    true,
			Tag:         "from",
			Label:       "From",
			Description: "Address of the account that will send the transferred tokens",
			DataType:    "->wallet",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		amountStr, _ := req["amount"].(string)
		recipientWallet, _ := req["to"].(assets.Key)
		senderWallet, _ := req["from"].(assets.Key)

		// check if amount is valid integer
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return nil, errors.WrapError(err, "invalid amount, send a positive integer")
		}

		if amount <= 0 {
			return nil, errors.WrapError(nil, "invalid amount, send a positive integer")
		}

		// get current amount
		recipientWalletAsset, err := recipientWallet.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting wallet asset")
		}

		recipientCurrentAmount, nerr := strconv.Atoi(recipientWalletAsset.GetProp("goTokenBalance").(string))
		if nerr != nil {
			return nil, errors.WrapError(err, "error converting current amount to integer")
		}

		recipientWalletMap := map[string]interface{}{
			"goTokenBalance": strconv.Itoa(recipientCurrentAmount + amount),
		}

		recipientWalletInLedger, err := recipientWallet.Update(stub, recipientWalletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error putting wallet asset")
		}

		// get current amount
		senderWalletAsset, err := senderWallet.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting wallet asset")
		}

		senderCurrentAmount, nerr := strconv.Atoi(senderWalletAsset.GetProp("goTokenBalance").(string))
		if nerr != nil {
			return nil, errors.WrapError(err, "error converting current amount to integer")
		}

		if senderCurrentAmount < amount {
			return nil, errors.WrapError(nil, "insufficient funds")
		}

		senderWalletMap := map[string]interface{}{
			"goTokenBalance": strconv.Itoa(senderCurrentAmount - amount),
		}

		senderWalletInLedger, err := senderWallet.Update(stub, senderWalletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error putting wallet asset")
		}

		response := map[string]interface{}{
			"recipientWallet": recipientWalletInLedger,
			"senderWallet":    senderWalletInLedger,
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling wallet asset")
		}

		return responseJSON, nil
	},
}
