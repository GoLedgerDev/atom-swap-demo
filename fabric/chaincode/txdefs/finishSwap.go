package txdefs

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"golang.org/x/crypto/sha3"
)

// POST Method
var FinishSwap = tx.Transaction{
	Tag:         "finishSwap",
	Label:       "Finish Swap",
	Description: "Finish a swap by sending the secret",
	Method:      "POST",

	Args: []tx.Argument{
		{
			Required: true,
			Tag:      "swap",
			Label:    "Swap to finish",
			DataType: "->hashTimeLock",
		},
		{
			Required: true,
			Tag:      "secret",
			Label:    "Secret",
			DataType: "string",
		},
	},

	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error

		swap, _ := req["swap"].(assets.Key)
		secret, _ := req["secret"].(string)

		// get swap asset
		swapAsset, err := swap.Get(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting swap asset")
		}

		swapAmountStr, _ := swapAsset.GetProp("amount").(string)
		swapAmount, _ := strconv.Atoi(swapAmountStr)

		// check if swap is already finished
		if swapAsset.GetProp("secret") != nil {
			return nil, errors.WrapError(nil, "swap is already finished")
		}

		// check if secret is correct
		hashlock, _ := swapAsset.GetProp("hashlock").(string)
		if hashlock != string(keccak256([]byte(secret))) {
			return nil, errors.WrapError(nil, "incorrect secret")
		}

		// check if timelock is expired
		timelock, _ := swapAsset.GetProp("timelock").(time.Time)

		nowTimestamp, err := stub.Stub.GetTxTimestamp()
		if err != nil {
			return nil, errors.WrapError(err, "error getting current timestamp")
		}

		now := nowTimestamp.AsTime()
		if now.After(timelock) {
			return nil, errors.WrapError(nil, "swap is expired")
		}

		// get fromWallet asset
		fw, _ := swapAsset.GetProp("fromWallet").(map[string]interface{})
		fromWallet, err := assets.NewKey(fw)
		if err != nil {
			return nil, errors.WrapError(err, "error creating fromWallet key")
		}
		fromWalletAsset, err := fromWallet.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting fromWallet asset")
		}

		blockedBalanceStr, _ := fromWalletAsset["blockedBalance"].(string)
		blockedBalance, _ := strconv.Atoi(blockedBalanceStr)

		newBlockedBalance := blockedBalance - swapAmount

		if newBlockedBalance < 0 {
			return nil, errors.WrapError(nil, "insufficient blocked funds")
		}

		fromWalletAsset["blockedBalance"] = strconv.Itoa(newBlockedBalance)

		_, err = fromWallet.Update(stub, fromWalletAsset)
		if err != nil {
			return nil, errors.WrapError(err, "error updating fromWallet asset")
		}

		// get toWallet asset
		tw, _ := swapAsset.GetProp("toWallet").(map[string]interface{})
		toWallet, err := assets.NewKey(tw)
		if err != nil {
			return nil, errors.WrapError(err, "error creating toWallet key")
		}
		toWalletAsset, err := toWallet.GetMap(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error getting toWallet asset")
		}

		balanceStr, _ := toWalletAsset["goTokenBalance"].(string)
		balance, _ := strconv.Atoi(balanceStr)

		newBalance := balance + swapAmount

		toWalletAsset["goTokenBalance"] = strconv.Itoa(newBalance)

		_, err = toWallet.Update(stub, toWalletAsset)
		if err != nil {
			return nil, errors.WrapError(err, "error updating toWallet asset")
		}

		updatedSwap, err := swap.Update(stub, map[string]interface{}{
			"secret": secret,
		})

		if err != nil {
			return nil, errors.WrapError(err, "error updating swap asset")
		}

		swapJSON, err := json.Marshal(updatedSwap)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling swap asset")
		}

		return swapJSON, nil
	},
}

func keccak256(data []byte) string {
	var hash [32]byte
	keccak := sha3.NewLegacyKeccak256()
	_, _ = keccak.Write(data)
	keccak.Sum(hash[:0])
	return fmt.Sprintf("%x", hash)
}
