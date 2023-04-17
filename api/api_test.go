package api

import (
	"btc_chain/rpc"
	"encoding/json"
	"testing"
)

func TestGetBlockInfoByBlockHash(t *testing.T) {
	client := rpc.NewBTCRPCHttpClient(
		"127.0.0.1:8332",
		"BtcMySlef",
		"Ploxs1",
	)
	bestBlockHash, err := client.Client.GetBestBlockHash()
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(bestBlockHash.String())
		//get the block hash,and get the block info
		blockInfo, err := GetBlockInfoByBlockHash(client.Client, bestBlockHash.String())
		if err != nil {
			t.Log("GetBlockInfoByBlockHash err:", err.Error())
		} else {
			bytes, _ := json.Marshal(blockInfo) //transfer the data to json
			t.Log(string(bytes))
		}
	}
}

func TestGetBlockCount(t *testing.T) {
	client := rpc.NewBTCRPCHttpClient(
		"127.0.0.1:8332",
		"BtcMySlef",
		"Ploxs1",
	)
	blockHeight, err := GetBlockCount(client.Client)
	if err != nil {
		t.Log("GetBlockCount err:", err.Error())
	} else {
		t.Log("BlockHeight======>", blockHeight)
	}
}

func TestGetBlockHashByBlockHeight(t *testing.T) {
	client := rpc.NewBTCRPCHttpClient(
		"127.0.0.1:8332",
		"BtcMySlef",
		"Ploxs1",
	)
	// get the height of 1's hash
	blockHash, err := GetBlockHashByBlockHeight(client.Client, 1)
	if err != nil {
		t.Log("GetBlockCount err:", err.Error())
	} else {
		t.Log("blockHash=====>", blockHash.String())
	}
}
