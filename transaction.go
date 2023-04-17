package btc_bl

import (
	"btc_chain/model"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"io"
	"net/http"
	"strings"
)

const GET = "GET"
const POST = "POST"

const (
	// blockCypher's token myself
	blockCypherAccessToken = "20056d07b28f458e90916d5209eecc87"
)

type UTXORet struct {
	Balance      int64 `json:"balance"`
	FinalBalance int64 `json:"final_balance"`
	Txrefs       []struct {
		TxHash      string `json:"tx_hash"`
		BlockHeight int64  `json:"block_height"`
		TxInputN    int64  `json:"tx_input_n"`
		TxOutputN   int64  `json:"tx_output_n"`
		BtcValue    int64  `json:"value"`
		Spent       bool   `json:"spent"`
	} `json:"txrefs"`
}

type SendTxRet struct {
	Tx struct {
		Hash string `json:"hash"`
	} `json:"tx"`
}

func getUTXOListFromBlockCypherAPI(address, netType string) (*UTXORet, error) {
	number := 1000 // limit get the max number 1000 UTXO from this address
	url := fmt.Sprintf("http://api.blockcypher.com/v1/%s/addrs/%s?unspentOnly=true&limit=%d"+
		"&includeScript=false&includeConfidence=false", netType, address, number) //combine the net chrome blockCypher API
	url = url + "&" + blockCypherAccessToken //take the auth token
	req, err := http.NewRequest(GET, url, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, nil
	}
	//analyze the data to the struct
	utxoList := UTXORet{}

	if err := json.Unmarshal(data, &utxoList); err != nil {
		return nil, err
	}
	return &utxoList, nil
}

func sendRawTransactionHexToNode_BlockCypherAPI(txHex, netType string) (string, error) {
	url := fmt.Sprintf("http://api.blockcypher.com/v1/%s/txs/push", netType)
	url = url + "?" + blockCypherAccessToken
	//{"tx":$TXHEX}
	jsonStr := fmt.Sprintf("{\"tx\":\"%s\"}", txHex) //constructor
	req, err := http.NewRequest(POST, url, strings.NewReader(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return "", err
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	fmt.Println("the result for request trans:", string(data))
	ret := SendTxRet{}
	if err := json.Unmarshal(data, &ret); err != nil {
		return "", err
	}
	return ret.Tx.Hash, nil
}

type OpReturnDataObj struct {
	Data string //set the data the opReturn
}

// send the test net btc transaction, 1 is stand for 0.000000001btc
func SendTestNet_BTCNormalTransaction(senderPrivateKey, toAddress string, value int64, opReturn *OpReturnDataObj) error {
	targetTransactionValue := btcutil.Amount(value)
	blockCypherApiTestNet := "btc/test3"
	// according to the sender's priKey to get info,ep:address
	wallet := CreateWalletFromPrivateKey(senderPrivateKey, &chaincfg.TestNet3Params)
	if wallet == nil {
		return errors.New("invalid private key") //recovery wallet failed
	}
	//prepare the utxo
	senderAddress := wallet.GetBtcAddress(true)
	utxoList, err := getUTXOListFromBlockCypherAPI(senderAddress, blockCypherApiTestNet)
	if err != nil {
		return err
	}
	//according to the transaction value choose the input data
	tx := wire.NewMsgTx(wire.TxVersion) //defined a transaction instance
	var (
		totalUTXOValue btcutil.Amount
		changeValue    btcutil.Amount
	)

	//SpendSize is BTC advise,be used count the gas fee
	SpendSize := 1 + 73 + 1 + 33
	for _, utxo := range utxoList.Txrefs {
		totalUTXOValue += btcutil.Amount(utxo.BtcValue) //statistics all the useful utxo
		hash := &chainhash.Hash{}
		if err := chainhash.Decode(hash, utxo.TxHash); err != nil {
			panic(fmt.Errorf("constractor hash failed:%s", err.Error()))
		}

		//use the pre transaction hash build the data input
		preUTXO := wire.OutPoint{Hash: *hash, Index: uint32(utxo.TxOutputN)}
		oneInput := wire.NewTxIn(&preUTXO, nil, nil)
		tx.AddTxIn(oneInput)

		//according to the data count the fee
		txSize := tx.SerializeSize() + SpendSize*len(tx.TxIn)
		reqFee := btcutil.Amount(txSize * 10)

		// the UTXO sub the hands fee and compare the TO aim
		if totalUTXOValue-reqFee < targetTransactionValue {
			// has not reach the value you want,while continue
			continue
		}

		// give change to self
		changeValue = totalUTXOValue - targetTransactionValue - reqFee
		break //reach the value,break,need not to add the utxo
	}

	//constructor the transaction output
	//come true the common wallet transaction
	toPubKeyHash := getAddressPubkeyHash(toAddress)
	if toPubKeyHash == nil {
		return errors.New("invalid receiver address") //nit legal wallet
	}
	toAddressPubKeyHashObj, err := btcutil.NewAddressPubKeyHash(toPubKeyHash, &chaincfg.TestNet3Params)
	if err != nil {
		return err
	}
	// toAddressLockScript is a lock script
	toAddressLockScript, err := txscript.PayToAddrScript(toAddressPubKeyHashObj)
	if err != nil {
		return err
	}
	//receiver is the inner (to) output
	receiverOutput := &wire.TxOut{PkScript: toAddressLockScript, Value: int64(targetTransactionValue)}
	tx.AddTxOut(receiverOutput) //add to the test struct
	var senderAddressLockScript []byte
	if changeValue > 0 { //if >0,you must give change to yourself
		//first you must to account the value senderAddressLockScript
		senderPubKeyHash := getAddressPubkeyHash(senderAddress)
		senderAddressPubKeyHashObj, err := btcutil.NewAddressPubKeyHash(senderPubKeyHash, &chaincfg.TestNet3Params)
		if err != nil {
			return err
		}
		//is a lock script
		senderAddressLockScript, err := txscript.PayToAddrScript(senderAddressPubKeyHashObj)
		if err != nil {
			return err
		}
		//give change
		senderOutput := &wire.TxOut{PkScript: senderAddressLockScript, Value: int64(changeValue)}
		//add to the transaction struct
		tx.AddTxOut(senderOutput)
	}
	btcecPrivatekey := (btcec.PrivateKey)(wallet.PrivateKey)
	txInsize := len(tx.TxIn)
	for i := 0; i < txInsize; i++ {
		sigScript, err := txscript.SignatureScript( //the script signature create
			tx,
			i,
			senderAddressLockScript,
			txscript.SigHashAll,
			&btcecPrivatekey,
			true)
		if err != nil {
			return err
		}
		tx.TxIn[i].SignatureScript = sigScript
	}

	//send Transaction
	//first count the transaction's hash,then send it to the transaction node,the data is hash
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return err
	}
	txHex := hex.EncodeToString(buf.Bytes())
	//send the hash data to the transaction
	txHash, err := sendRawTransactionHexToNode_BlockCypherAPI(txHex, blockCypherApiTestNet)
	if err != nil {
		return err
	}
	fmt.Println("the transaction hash:", txHash)

	if opReturn != nil {
		nullDataScript, err := txscript.NullDataScript([]byte(opReturn.Data))
		if err != nil {
			return err
		}
		opreturnOutput := &wire.TxOut{PkScript: nullDataScript, Value: 0}
		tx.AddTxOut(opreturnOutput)
	}

	return nil
}

// get the address's pubKeyHash
func getAddressPubkeyHash(address string) []byte {
	pubKeyHash := model.Base58Encode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	return pubKeyHash
}
