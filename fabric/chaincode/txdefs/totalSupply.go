package txdefs

import (
	"encoding/json"
	"net/http"

	"github.com/goledgerdev/cc-tools/assets"
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
	"github.com/goledgerdev/token-cc/chaincode/assettypes"
)

// GET Method
var TotalSupply = tx.Transaction{
	Tag:         "totalSupply",
	Label:       "Get Total Supply",
	Description: "TotalSupply returns the total token supply",
	Method:      "GET",

	Args: []tx.Argument{
		{},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {
		var err error

		totalSupplyKey, err := assets.NewKey(map[string]interface{}{
			"@assetType":     "totalSupply",
			"totalSupplyKey": assettypes.TotalSupplyKey,
		})
		if err != nil {
			return nil, errors.WrapError(err, "error creating total supply key")
		}

		totalSupplyExists, err := totalSupplyKey.ExistsInLedger(stub)
		if err != nil {
			return nil, errors.WrapError(err, "error checking if total supply asset exists")
		}

		response := map[string]interface{}{
			"totalSupply": "0",
		}

		if totalSupplyExists {
			totalSupplyObj, err := totalSupplyKey.Get(stub)
			if err != nil {
				return nil, errors.WrapError(err, "error getting total supply asset")
			}

			totalSupply, ok := totalSupplyObj.GetProp("supply").(string)
			if !ok {
				return nil, errors.NewCCError("error getting total supply", http.StatusInternalServerError)
			}

			response["totalSupply"] = totalSupply
		}

		responseJSON, err := json.Marshal(response)
		if err != nil {
			return nil, errors.WrapError(err, "error marshalling response")
		}

		return responseJSON, nil
	},
}
