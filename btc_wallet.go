package btc_bl

import (
	"btc_chain/model"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

// the BTC wallet struct
type BtcWallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  *btcec.PublicKey
	chainType  *chaincfg.Params
}

// according to the type of chain to create wallet
func CreateBtcWallet(chainType *chaincfg.Params) *BtcWallet {
	private, public := newKeyPair() //use elliptic curve create the privateKey and publicKey
	wallet := BtcWallet{PrivateKey: private, PublicKey: public}
	wallet.chainType = chainType
	return &wallet
}

func newKeyPair() (ecdsa.PrivateKey, *btcec.PublicKey) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	_, pubKey := btcec.PrivKeyFromBytes(btcec.S256(), private.D.Bytes())
	return *private, pubKey
}

// import the wtf privateKey ,recover wallet
func CreateWalletFromPrivateKey(wifPrivatekey string, chainType *chaincfg.Params) *BtcWallet {
	wif, err := btcutil.DecodeWIF(wifPrivatekey)
	if err != nil {
		panic(err)
	}
	privKeyBytes := wif.PrivKey.Serialize()
	priKey, publicKey := btcec.PrivKeyFromBytes(btcec.S256(), privKeyBytes)
	return &BtcWallet{
		PrivateKey: ecdsa.PrivateKey(*priKey),
		chainType:  chainType,
		PublicKey:  publicKey,
	}
}

func (w *BtcWallet) GetWIFPrivateKey(compress bool) string {
	if w.PrivateKey.D == nil {
		return ""
	}
	var combine = []byte{}
	//0x80 is WIF version
	if compress {
		//zip version
		combine = append([]byte{0x80}, w.PrivateKey.D.Bytes()...)
		combine = append(combine, 0x01)
	} else {
		combine = append([]byte{0x80}, w.PrivateKey.D.Bytes()...)
	}
	checkCodeBytes := doubleSha256F(combine)
	combine = append(combine, checkCodeBytes[0:4]...)
	//baseCoding
	return string(model.Base58Encode(combine))
}

func doubleSha256F(payload []byte) []byte {
	sha256 := sha256.New()
	sha256.Reset()
	sha256.Write(payload)
	hash1 := sha256.Sum(nil)
	sha256.Reset()
	sha256.Write(hash1)
	return sha256.Sum(nil)
}

func (w *BtcWallet) GetPubKeyHexStr(compress bool) string {
	//zip for
	if compress {
		return hex.EncodeToString(w.PublicKey.SerializeCompressed())
	}
	//comp for
	return hex.EncodeToString(w.PublicKey.SerializeUncompressed())
}

func (w *BtcWallet) GetBtcAddress(compress bool) string {
	var buf []byte
	if compress {
		//zip
		buf = w.PublicKey.SerializeCompressed()
	} else {
		buf = w.PublicKey.SerializeUncompressed()
	}
	if buf == nil {
		return ""
	}

	//Hash160 inner--->Sha256 and RIPEMD160
	pubKeyHash := btcutil.Hash160(buf)
	addr, err := btcutil.NewAddressPubKeyHash(pubKeyHash, w.chainType)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	//addr.String checkCode&&Base58
	return addr.String()
}

// the struct for multisig wallet
type MultiWallet struct {
	Address      string `json:"address"`
	RedeemScript string `json:"redeemScript"`
}

// @Param n: you should use n PrivKey to unlock
// @Param pubKeyStrs: PubKeys array
// @Param chainType: is the type of the net
func CreateMultiWallet(n int, pubKeyStrs []string, chainType *chaincfg.Params) (*MultiWallet, error) {
	pubKeyList := make([]*btcutil.AddressPubKey, len(pubKeyStrs))
	//add all the pubKeys
	for index, pubKeyItem := range pubKeyStrs {
		//Use Decode transfer the pubKey to the address
		addressObj, err := btcutil.DecodeAddress(pubKeyItem, chainType)
		if err != nil {
			return nil, fmt.Errorf("invalid pubKey %s :decode failed error %s", pubKeyItem, err.Error())
		}
		// judge the chain's type equal the chainType
		if !addressObj.IsForNet(chainType) {
			return nil, fmt.Errorf("invalid pubKey %s:not match chain type %s", pubKeyItem, err.Error())
		}
		switch pubKey := addressObj.(type) {
		case *btcutil.AddressPubKey:
			pubKeyList[index] = pubKey
			continue
		default:
			// discovery not match the data for address,ep:the param is not the pubKey
			return nil, fmt.Errorf("address contains invalid item")
		}
	}
	script, err := txscript.MultiSigScript(pubKeyList, n)
	if err != nil {
		return nil, fmt.Errorf("failed to parse scirpt %s", err.Error())
	}
	address, err := btcutil.NewAddressScriptHash(script, chainType)
	if err != nil {
		return nil, err
	}
	return &MultiWallet{
		Address: address.EncodeAddress(),
		//use the hex
		RedeemScript: hex.EncodeToString(script),
	}, nil
}
