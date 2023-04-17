package api

import (
	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

// according to the transaction hash to get the transaction data
func GetTransactionInfoByTxHash(client *rpcclient.Client, txHash string) (*btcjson.TxRawResult, error) {
	hash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, err
	}
	return client.GetRawTransactionVerbose(hash)
}
