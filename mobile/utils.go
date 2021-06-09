package mobile

// https://github.com/algorand/go-algorand-sdk/compare/a140151ac15136f234b9cf094e58cb69cb37e2c0...MobileCompatible

import (
	"bytes"
	"crypto/sha512"
	"encoding/base32"
	"encoding/base64"
	"fmt"
	"math"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/types"
)

type Uint64 struct {
	Upper int64
	Lower int64
}

func MakeUint64(value uint64) Uint64 {
	return Uint64{
		Upper: int64(value >> 32),
		Lower: int64(math.MaxUint32 & value),
	}
}

func (i Uint64) Extract() (value uint64, err error) {
	if i.Upper < 0 || i.Upper > math.MaxUint32 {
		err = fmt.Errorf("Upper value of Uint64 not in correct range. Expected value between 0 and %d, got %d", math.MaxUint32, i.Upper)
		return
	}

	if i.Lower < 0 || i.Lower > math.MaxUint32 {
		err = fmt.Errorf("Lower value of Uint64 not in correct range. Expected value between 0 and %d, got %d", math.MaxUint32, i.Lower)
		return
	}

	value = uint64(i.Upper)<<32 | uint64(i.Lower)

	return
}

func IsValidAddress(addr string) bool {
	_, err := types.DecodeAddress(addr)
	return err == nil
}

// PendingTransactionsResponse is returned by PendingTransactions and by Txid
type PendingTransactionsResponse = struct {
	TopTransactions   []types.SignedTxn `codec:"top-transactions"`
	TotalTransactions uint64            `codec:"total-transactions"`
}

// TODO: fix this method, since it returns []string, which is not supported by gomobile
// do we even need this?
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
		res = append(res, txIDFromTransaction(txn.Txn))
	}

	return res, nil
}

// rawTransactionBytesToSign returns the byte form of the tx that we actually sign
// and compute txID from.
func rawTransactionBytesToSign(tx types.Transaction) []byte {
	// Encode the transaction as msgpack
	encodedTx := msgpack.Encode(tx)

	// Prepend the hashable prefix
	msgParts := [][]byte{[]byte("TX"), encodedTx}
	return bytes.Join(msgParts, nil)
}

// transactionID is the unique identifier for a Transaction in progress
func transactionID(tx types.Transaction) (txid []byte) {
	toBeSigned := rawTransactionBytesToSign(tx)
	txid32 := sha512.Sum512_256(toBeSigned)
	txid = txid32[:]
	return
}

// txIDFromTransaction is a convenience function for generating txID from txn
func txIDFromTransaction(tx types.Transaction) (txid string) {
	txidBytes := transactionID(tx)
	txid = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(txidBytes[:])
	return
}
