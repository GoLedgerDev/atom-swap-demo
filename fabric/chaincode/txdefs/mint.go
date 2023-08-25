package txdefs

import (
	"github.com/goledgerdev/cc-tools/errors"
	sw "github.com/goledgerdev/cc-tools/stubwrapper"
	tx "github.com/goledgerdev/cc-tools/transactions"
)

// POST Method
var Mint = tx.Transaction{
	Tag:         "mint",
	Label:       "Mint",
	Description: "Mint creates new tokens and adds them to minter's account balance",
	Method:      "POST",

	Args: []tx.Argument{
		{},
	},
	Routine: func(stub *sw.StubWrapper, req map[string]interface{}) ([]byte, errors.ICCError) {

		return []byte{}, nil
	},
}
