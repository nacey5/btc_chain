package btc_bl

import (
	"github.com/btcsuite/btcd/chaincfg"
	"testing"
)

func TestCreateBtcWallet(t *testing.T) {
	//create the main chain wallet
	btcWallet := CreateBtcWallet(&chaincfg.MainNetParams)
	t.Log("zip for:",btcWallet.GetWIFPrivateKey(true)) //print the zip wif privateKey
	t.Log("comp for:",btcWallet.GetWIFPrivateKey(false))
}