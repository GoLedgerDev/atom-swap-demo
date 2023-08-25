package main

import (
	"github.com/goledgerdev/cc-tools/events"
	"github.com/goledgerdev/token-cc/chaincode/eventtypes"
)

var eventTypeList = []events.Event{
	eventtypes.CreateLibraryLog,
}
