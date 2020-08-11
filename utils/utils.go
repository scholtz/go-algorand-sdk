package utils

import (
	"encoding/base64"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

func IsValidAddress(addr string) bool {
	_, err := types.DecodeAddress(addr)
	return err == nil
}

// PendingTransactionsResponse is returned by PendingTransactions and by Txid
type PendingTransactionsResponse = struct {
	TopTransactions   []types.SignedTxn `codec:"top-transactions"`
	TotalTransactions uint64            `codec:"total-transactions"`
}

func PendingResultsToTXID(response string) ([]string, error) {
	sDec, err := base64.StdEncoding.DecodeString(response)

	if err != nil {
		return nil, err
	}

	var result PendingTransactionsResponse
	err = msgpack.Decode(sDec, &result)

	if err != nil {
		return nil, err
	}

	var res []string
	for _, txn := range result.TopTransactions {
		res = append(res, types.TxIDFromTransaction(txn.Txn))
	}

	return res, nil
}
