package xc_share_file_client

import (
	"crypto/rsa"
	"encoding/base32"
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
	privateKeyBS := crypt_tools.RSAPrivateKeyToDer(privateKey)
	clientPrivateKey := base32.StdEncoding.EncodeToString(privateKeyBS)

	client := xchg_connections.NewClientConnection(xchg_network.NewNetworkDefault(), publicAddress, clientPrivateKey, password, nil)
	var bs []byte
	bs, err = client.Call("get-file-content", nil)
	if err != nil {
		fmt.Println("ERROR:", err)
		return
	}
	fmt.Println("RECEIVED:", string(bs))
}
