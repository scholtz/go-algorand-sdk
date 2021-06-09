package mobile

import (
	"github.com/algorand/go-algorand-sdk/encoding/json"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

// TransactionMsgpackToJson converts a msgpack-encoded Transaction to a
// json-encoded Transaction
func TransactionMsgpackToJson(msgpTxn []byte) (jsonTxn []byte, err error) {
	var txn types.Transaction
	err = msgpack.Decode(msgpTxn, &txn)
	if err == nil {
		jsonTxn = json.Encode(txn)
	}
	return
}

// TransactionJsonToMsgpack converts a json-encoded Transaction to a
// msgpack-encoded Transaction
func TransactionJsonToMsgpack(jsonTxn []byte) (msgpackTxn []byte, err error) {
	var txn types.Transaction
	err = json.Decode(jsonTxn, &txn)
	if err == nil {
		msgpackTxn = msgpack.Encode(txn)
	}
	return
}
