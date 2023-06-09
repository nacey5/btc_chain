package btc_bl

import (
	"encoding/hex"
	"testing"
)

func TestGetUTXO(t *testing.T) {
	ret, err := getUTXOListFromBlockCypherAPI("364zJb64c1hYLQ2NTtk68gJYpochH2Xe6w", "btc/test3")
	t.Log(err)
	t.Log(ret)
}

func TestSendTestNet_BTCNormalTransaction(t *testing.T) {
	err := SendTestNet_BTCNormalTransaction(
		"KzZkYh62v6xq2SdMaYbuR6yhbbav1Pq9cXGU6M8Ci8m6J6qc23r3",
		"1KFHE7w8BhaENAswwryaoccDb6qcT6DbYY",
		2000, &OpReturnDataObj{
			Data: "This is my transaction info",
		})
	t.Log(err)
}

func TestDecodeOpeReturnData(t *testing.T) {
	hexData := "KzZkYh62v6xq2SdMaYbuR6yhbbav1Pq9cXGU6M8Ci8m6J6qc23r3"
	originDataBytes, _ := hex.DecodeString(hexData)
	t.Log(string(originDataBytes))
}
