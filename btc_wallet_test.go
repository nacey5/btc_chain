package btc_bl

import (
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestCreateBtcWallet(t *testing.T) {
	//create the main chain wallet
	btcWallet := CreateBtcWallet(&chaincfg.MainNetParams)
	t.Log("zip for:", btcWallet.GetWIFPrivateKey(true)) //print the zip wif privateKey
	t.Log("comp for:", btcWallet.GetWIFPrivateKey(false))
	t.Log(">>>>>>>>>>>>>>>>>>>PUB")
	t.Log("zip for:", btcWallet.GetPubKeyHexStr(true))
	t.Log("comp for:", btcWallet.GetPubKeyHexStr(false))
}

func TestCreateMultiWallet(t *testing.T) {
	pubKeys := []string{
		"04fa9d1d45681f24416a7277d756c8e2aa0cf07a35a35d1f27e013634c72cf5a365172cdc2cd209c0b0057e18ba831a5f5564296f24ae130e4fb551f26f2a4a121",
		"04f0f44dda6bfecb024fc071d4d7ff237b281297bd8156807c71f0699147aa930aba5d73db8575111762b404ab57d258126cb12bb8ec62b6217323d3e591e4a27e",
		"047713603eb1db0656b9a0917a88118f9bc1917cf06ee7a49e76e4f3b0059fb363791a2ef85b82de42a51c2457e03684ca209e72ba7c538b49f19f41d537aa62a7",
	}
	mutil, err := CreateMultiWallet(2, pubKeys, &chaincfg.MainNetParams)
	t.Log(err)
	if err == nil {
		t.Log(mutil.Address)
	}
}
