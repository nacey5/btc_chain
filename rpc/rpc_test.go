package rpc

import "testing"

func TestRPCClient(t *testing.T) {
	client := NewBTCRPCHttpClient(
		"127.0.0.1:8332",
		"BtcMySlef",
		"Ploxs1",
	)
	bestBlockHash, err := client.Client.GetBestBlockHash()
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(bestBlockHash.String())
	}
}
