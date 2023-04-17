package rpc

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
)

type BTCRPCClient struct {
	NodeUrl string
	Client  *rpcclient.Client //stand for the rpc client
}

func NewBTCRPCHttpClient(nodeUrl, user, password string) *BTCRPCClient {
	connCfg := &rpcclient.ConnConfig{
		Host:         nodeUrl,
		User:         user,
		Pass:         password,
		HTTPPostMode: true,
		DisableTLS:   true,
		//DisableTLS:
		// if RPC service enable the https ,advise use tls always
		// set the value is false, if not do this,the user and pass will send as the orign value
		// if not use http,that set the value is true
	}
	rpcClient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		//init failed ,end the app,and show the error msg to the console
		errInfo := fmt.Errorf("init rpc client failed %s", err.Error())
		panic(errInfo)
	}
	return &BTCRPCClient{
		NodeUrl: nodeUrl,
		Client:  rpcClient,
	}
}

func NewBTCRPCSocketClient(nodeUrl, user, password string) *BTCRPCClient {
	connCfg := &rpcclient.ConnConfig{
		Host:         nodeUrl,
		Endpoint:     "ws",
		User:         user,
		Pass:         password,
		Certificates: nil,
		//Certificates: if enable https.must config the Certificates
	}
	handlers := rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txs []*btcutil.Tx) {
			//in this func ,can callback
		},
	}
	rpcClient, err := rpcclient.New(connCfg, &handlers)
	if err != nil {
		errInfo := fmt.Errorf("init rpc client failed %s", err.Error())
		panic(errInfo)
	}
	return &BTCRPCClient{
		NodeUrl: nodeUrl,
		Client:  rpcClient,
	}
}
