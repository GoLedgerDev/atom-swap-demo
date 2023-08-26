package txdefs

import (
	"encoding/json"
	"strconv"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/goledgerdev/token-cc/chaincode/assettypes"
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
			DataType:    "->wallet",
		},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		amountStr, _ := req["amount"].(string)
		wallet, _ := req["to"].(assets.Key)

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

		walletMap := map[string]interface{}{
			"goTokenBalance": strconv.Itoa(currentAmount + amount),
		}

		walletInLedger, err := wallet.Update(stub, walletMap)
		if err != nil {
			return nil, errors.WrapError(err, "error putting wallet asset")
		}

		walletJSON, err := json.Marshal(walletInLedger)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling wallet asset")
		}

		// add amount to total supply
		err = AddAmountToTotalSupply(stub, amount)
		if err != nil {
			return nil, errors.WrapError(err, "error adding amount to total supply")
		}

		return walletJSON, nil
	},
}

func AddAmountToTotalSupply(stub *sw.StubWrapper, amountToAdd int) error {
	totalSupplyKey, err := assets.NewKey(map[string]interface{}{
		"@assetType":     "totalSupply",
		"totalSupplyKey": assettypes.TotalSupplyKey,
	})
	if err != nil {
		return errors.WrapError(err, "error creating total supply key")
	}

	totalSupplyExists, err := totalSupplyKey.ExistsInLedger(stub)
	if err != nil {
		return errors.WrapError(err, "error checking if total supply asset exists")
	}

	totalSupplyAsset, err := assets.NewAsset(map[string]interface{}{
		"@assetType":     "totalSupply",
		"totalSupplyKey": assettypes.TotalSupplyKey,
	})

	if !totalSupplyExists {
		err := totalSupplyAsset.SetProp("supply", strconv.Itoa(amountToAdd))
		if err != nil {
			return errors.WrapError(err, "error setting total supply asset")
		}
	} else {
		var err error
		totalSupplyObj, err := totalSupplyKey.Get(stub)
		if err != nil {
			return errors.WrapError(err, "error getting total supply asset")
		}

		currentTotalSupply, err := strconv.Atoi(totalSupplyObj.GetProp("supply").(string))
		if err != nil {
			return errors.WrapError(err, "error converting current supply to integer")
		}

		err = totalSupplyAsset.SetProp("supply", strconv.Itoa(currentTotalSupply+amountToAdd))
		if err != nil {
			return errors.WrapError(err, "error setting total supply asset")
		}
	}

	_, err = totalSupplyAsset.Put(stub)
	if err != nil {
		return errors.WrapError(err, "error putting total supply asset")
	}

	return nil
}
