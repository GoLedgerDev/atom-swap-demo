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
			DataType:    "string",
		},
		{
			Required:    true,
			Tag:         "from",
			Label:       "From",
			Description: "Address of the account that will send the transferred tokens",
			DataType:    "string",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		amountStr, _ := req["amount"].(string)
		to, _ := req["to"].(string)
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
		_, err = utils.CheckPublicKey(to)
		if err != nil {
			return nil, errors.WrapError(nil, "invalid recipient address")
		}

		// check if address is valid
		_, err = utils.CheckPublicKey(from)
		if err != nil {
			return nil, errors.WrapError(nil, "invalid sender address")
		}

		// retrieve recipient recipientWallet
		recipientWallet, err := assets.NewKey(map[string]interface{}{
			"@assetType": "wallet",
			"address":    to,
		})

		// check if wallet exists
		exists, err := recipientWallet.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error checking if wallet exists")
		}

		recipientCurrentAmount := 0
		if exists {
			// get current amount
			recipientWalletAsset, err := recipientWallet.Get(stub)
			if err != nil {
				return nil, err
			}

			c, nerr := strconv.Atoi(recipientWalletAsset.GetProp("goTokenBalance").(string))
			if nerr != nil {
				return nil, errors.WrapError(err, "error converting current amount to integer")
			}

			recipientCurrentAmount = c
		}

		recipientWalletMap := map[string]interface{}{
			"@assetType":     "wallet",
			"address":        to,
			"goTokenBalance": strconv.Itoa(recipientCurrentAmount + amount),
		}

		rw, err := assets.NewAsset(recipientWalletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error creating wallet asset")
		}

		recipientWalletInLedger, err := rw.Put(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error putting wallet asset")
		}

		// retrieve sender senderWallet
		senderWallet, err := assets.NewKey(map[string]interface{}{
			"@assetType": "wallet",
			"address":    from,
		})

		// check if wallet exists
		exists, err = senderWallet.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error checking if wallet exists")
		}

		senderCurrentAmount := 0
		if exists {
			// get current amount
			senderWalletAsset, err := senderWallet.Get(stub)
			if err != nil {
				return nil, err
			}

			c, nerr := strconv.Atoi(senderWalletAsset.GetProp("goTokenBalance").(string))
			if nerr != nil {
				return nil, errors.WrapError(err, "error converting current amount to integer")
			}

			senderCurrentAmount = c
		}

		if senderCurrentAmount < amount {
			return nil, errors.WrapError(nil, "insufficient funds")
		}

		senderWalletMap := map[string]interface{}{
			"@assetType":     "wallet",
			"address":        from,
			"goTokenBalance": strconv.Itoa(senderCurrentAmount - amount),
		}

		sw, err := assets.NewAsset(senderWalletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error creating wallet asset")
		}

		senderWalletInLedger, err := sw.Put(stub)
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
