package xc_share_file_client

import (
	"crypto/rsa"
	"fmt"

	"github.com/ipoluianov/gomisc/crypt_tools"
	"github.com/ipoluianov/xchg/xchg_connections"
	"github.com/ipoluianov/xchg/xchg_network"
)

func GetFile(publicAddress string, password string) {
	var err error
	var privateKey *rsa.PrivateKey
	privateKey, err = crypt_tools.GenerateRSAKey()
	if err != nil {
		return
	}
	clientPrivateKey := crypt_tools.RSAPrivateKeyToBase58(privateKey)

	client := xchg_connections.NewClientConnection(xchg_network.NewNetworkDefault(), publicAddress, clientPrivateKey, password)
	var bs []byte
	bs, err = client.Call("get-file-content", nil)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println("RECEIVED:", string(bs))
}
