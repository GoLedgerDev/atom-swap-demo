package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goledgerdev/ccapi/chaincode"
	"github.com/goledgerdev/ccapi/common"
)

func InvokeV1(c *gin.Context) {
	// Get channel information from request
	req := make(map[string]interface{})
	err := c.BindJSON(&req)
	if err != nil {
		common.Abort(c, http.StatusBadRequest, err)
		return
	}

	channelName := os.Getenv("CHANNEL")
	chaincodeName := os.Getenv("CCNAME")
	txName := c.Param("txname")

	var collections []string
	collectionsQuery := c.Query("@collections")
	if collectionsQuery != "" {
		collectionsByte, err := base64.StdEncoding.DecodeString(collectionsQuery)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the @collections query parameter must be a base64-encoded JSON array of strings",
			})
			return
		}

		err = json.Unmarshal(collectionsByte, &collections)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "the @collections query parameter must be a base64-encoded JSON array of strings",
			})
			return
		}
	} else {
		collectionsQuery := c.QueryArray("collections")
		if len(collectionsQuery) > 0 {
			collections = collectionsQuery
		} else {
			collections = []string{c.Query("collections")}
		}
	}

	transientMap := make(map[string]interface{})
	for key, value := range req {
		if key[0] == '~' {
			keyTrimmed := strings.TrimPrefix(key, "~")
			transientMap[keyTrimmed] = value
			delete(req, key)
		}
	}

	args, err := json.Marshal(req)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	transientMapByte, err := json.Marshal(transientMap)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	argList := [][]byte{}
	if args != nil {
		argList = append(argList, args)
	}

	res, status, err := chaincode.Invoke(channelName, chaincodeName, txName, argList, transientMapByte)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	var payload interface{}
	err = json.Unmarshal(res.Payload, &payload)
	if err != nil {
		common.Abort(c, http.StatusInternalServerError, err)
		return
	}

	common.Respond(c, payload, status, err)
}
