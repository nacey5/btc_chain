package api

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

func GetBlockInfoByBlockHash(client *rpcclient.Client, blockHash string) (*wire.MsgBlock, error) {
	hash, err := chainhash.NewHashFromStr(blockHash) //if the block hash is not hex ,will happen some error
	if err != nil {
		return nil, err
	}
	return client.GetBlock(hash)
}

// 获取链的最新区块高度
func GetBlockCount(client *rpcclient.Client) (int64, error) {
	return client.GetBlockCount()
}

// 根据区块高度获取区块哈希值
func GetBlockHashByBlockHeight(client *rpcclient.Client, height int64) (*chainhash.Hash, error) {
	return client.GetBlockHash(height)
}
