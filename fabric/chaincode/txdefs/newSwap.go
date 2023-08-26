package txdefs

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// POST Method
var NewSwap = tx.Transaction{
	Tag:         "newSwap",
	Label:       "New Swap",
	Description: "Create a new swap",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "id",
			Label:    "Lock ID",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "fromWallet",
			Label:    "From Wallet",
			DataType: "->wallet",
		},
		{
			Required: true,
			Tag:      "toWallet",
			Label:    "To Wallet",
			DataType: "->wallet",
		},
		{
			Required: true,
			Tag:      "amount",
			Label:    "Amount",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "hashlock",
			Label:    "Hash Lock",
			DataType: "string",
		},
		{
			Required: true,
			Tag:      "timelock",
			Label:    "Time Lock",
			DataType: "datetime",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error

		id, _ := req["id"].(string)
		fromWallet, _ := req["fromWallet"].(assets.Key)
		toWallet, _ := req["toWallet"].(assets.Key)
		amountStr, _ := req["amount"].(string)
		hashlock, _ := req["hashlock"].(string)
		timelock, _ := req["timelock"].(time.Time)

		// check if amount is valid integer
		amount, err := strconv.Atoi(amountStr)
		if err != nil {
			return nil, errors.WrapError(err, "invalid amount, send a positive integer")
		}

		if amount <= 0 {
			return nil, errors.WrapError(nil, "invalid amount, send a positive integer")
		}

		// get current amount from fromWallet
		walletAsset, err := fromWallet.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting wallet asset")
		}

		currentAmount, nerr := strconv.Atoi(walletAsset.GetProp("goTokenBalance").(string))
		if nerr != nil {
			return nil, errors.WrapError(err, "error converting current amount to integer")
		}

		if currentAmount < amount {
			return nil, errors.WrapError(err, "insufficient funds")
		}

		currentBlockedAmount, nerr := strconv.Atoi(walletAsset.GetProp("blockedBalance").(string))
		if nerr != nil {
			return nil, errors.WrapError(err, "error converting blocked amount to integer")
		}

		walletMap := map[string]interface{}{
			"goTokenBalance": strconv.Itoa(currentAmount - amount),
			"blockedBalance": strconv.Itoa(currentBlockedAmount + amount),
		}

		_, err = fromWallet.Update(stub, walletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error putting wallet asset")
		}

		htlAsset, err := assets.NewAsset(map[string]interface{}{
			"@assetType": "hashTimeLock",
			"id":         id,
			"fromWallet": fromWallet,
			"toWallet":   toWallet,
			"amount":     amountStr,
			"hashlock":   hashlock,
			"timelock":   timelock,
		})

		if err != nil {
			return nil, errors.WrapError(err, "failed to create wallet")
		}

		htlMap, err := htlAsset.PutNew(stub)
		if err != nil {
			return nil, errors.WrapError(err, "failed to create wallet")
		}

		htlJSON, err := json.Marshal(htlMap)
		if err != nil {
			return nil, errors.WrapError(err, "failed to create wallet")
		}

		return htlJSON, nil
	},
}
